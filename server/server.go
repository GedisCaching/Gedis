package redis

import (
	"errors"
	"github.com/GedisCaching/Gedis/storage"
	"strings"
	"sync"
)

// Server represents Redis-like server
type Server struct {
	// Data stores
	db *storage.Database

	// Configuration
	config *Config
}

// ServerManager manages multiple server instances by configuration
type ServerManager struct {
	servers map[Config]*Server

	// Mutex for thread safety when accessing the servers map
	mu sync.RWMutex
}

// Global server manager instance
var globalManager = &ServerManager{
	servers: make(map[Config]*Server),
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
		return server, nil
	}

	// Create a copy of the config to store in the server
	configCopy := config

	// Create new server with the config
	server := &Server{
		db:     storage.NewDatabase(),
		config: &configCopy,
	}

	// Store in map
	sm.servers[config] = server

	return server, nil
}
