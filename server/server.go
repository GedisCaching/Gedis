package redis

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/GedisCaching/Gedis/storage"
)

// Server represents Redis-like server
type Server struct {
	// Data stores
	db *storage.Database

	// Configuration
	config *Config

	// Last access time for LRU
	lastAccessed time.Time
}

// GetDB returns the database instance
func (s *Server) GetDB() *storage.Database {
	return s.db
}

// UpdateAccessTime updates the lastAccessed time of this server
// and its position in the LRU list
func (s *Server) UpdateAccessTime() {
	// Don't update the timestamp here, let the manager do it
	// Just notify the manager that this server was accessed
	globalManager.UpdateServerAccess(*s.config)
}

// ServerManager manages multiple server instances by configuration
type ServerManager struct {
	servers map[Config]*Server

	// LRU implementation
	capacity int
	lruList  []Config // List to track order for LRU (first item is least recently used)

	// Mutex for thread safety when accessing the servers map
	mu sync.RWMutex
}

// Global server manager instance
var globalManager = &ServerManager{
	servers:  make(map[Config]*Server),
	capacity: 100, // Default capacity, can be configured
	lruList:  make([]Config, 0),
}

// UpdateServerAccess updates a server's position in the LRU list
// This should be called whenever a server is accessed for any operation
func (sm *ServerManager) UpdateServerAccess(config Config) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if server, exists := sm.servers[config]; exists {
		// Update the server's last accessed time with current time
		now := time.Now()
		server.lastAccessed = now

		// Update its position in the LRU list
		sm.updateLRU(config)
	}
}

// evictIfNeeded removes least recently used servers if capacity is exceeded
func (sm *ServerManager) evictIfNeeded() {
	for len(sm.servers) > sm.capacity && len(sm.lruList) > 0 {
		// Remove the least recently used server (first in lruList)
		config := sm.lruList[0]
		sm.lruList = sm.lruList[1:] // Remove first element

		// Clean up resources for the server before deleting
		if server, exists := sm.servers[config]; exists {
			// Additional cleanup if needed
			_ = server // Avoid unused variable warning
		}

		delete(sm.servers, config)
	}
}

// updateLRU moves the accessed config to the end of the LRU list (most recently used)
func (sm *ServerManager) updateLRU(config Config) {
	// Find and remove the config from the current list
	idx := -1
	for i, cfg := range sm.lruList {
		if cfg == config {
			idx = i
			break
		}
	}

	if idx >= 0 {
		// Remove from current position
		sm.lruList = append(sm.lruList[:idx], sm.lruList[idx+1:]...)
	}

	// Add to the end of the list (most recently used position)
	sm.lruList = append(sm.lruList, config)
}

// NewServer creates a new server with default config if none provided
func NewServer(config *Config) (*Server, error) {
	if config == nil {
		config = DefaultConfig()
	}
	return NewServerWithConfig(config)
}

// NewServerWithConfig creates a new Redis-like server with custom configuration
// or returns an existing one if it already exists
func NewServerWithConfig(config *Config) (*Server, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// We need to dereference the pointer to use the Config as a map key
	return globalManager.GetOrCreateServer(*config)
}

func validateConfig(config *Config) error {
	if config.Address == "" {
		return errors.New("address cannot be empty")
	}

	// If address contains a port, validate it
	if strings.Contains(config.Address, ":") {
		parts := strings.Split(config.Address, ":")
		if len(parts) > 1 && parts[1] == "" {
			return errors.New("port cannot be empty if colon is present")
		}
	}

	return nil
}

// GetOrCreateServer returns an existing server for the given config or creates a new one
func (sm *ServerManager) GetOrCreateServer(config Config) (*Server, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if server already exists
	if server, exists := sm.servers[config]; exists {
		// Update last accessed time
		server.lastAccessed = time.Now()

		// Update LRU order
		sm.updateLRU(config)

		return server, nil
	}

	// Evict least recently used servers if we're at capacity
	sm.evictIfNeeded()

	// Create a copy of the config to store in the server
	configCopy := config

	// Create new server with the config
	server := &Server{
		db:           storage.NewDatabase(),
		config:       &configCopy,
		lastAccessed: time.Now(),
	}

	// Store in map
	sm.servers[config] = server

	// Add to LRU tracking
	sm.updateLRU(config)

	return server, nil
}

// GetServerCount returns the current number of servers in the manager
func (sm *ServerManager) GetServerCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.servers)
}
