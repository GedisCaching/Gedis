package storage

import (
	"sync"
	"time"
)

// DB represents the main database interface
type DB interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string) bool
	SetWithExpiry(key string, value interface{}, expiry time.Duration)
	Keys() []string
}

// Database represents "in-memory" Redis-like database
type Database struct {
	data    map[string]interface{}
	mu      sync.RWMutex // Mutex for concurrent access
	expires map[string]time.Time
}

// NewDatabase creates a new "in-memory" database
func NewDatabase() *Database {
	return &Database{
		data:    make(map[string]interface{}),
		expires: make(map[string]time.Time),
	}
}

// Set stores a key-value pair
func (db *Database) Set(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// Get retrieves a value by key
func (db *Database) Get(key string) (interface{}, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	value, exists := db.data[key]
	if !exists {
		return nil, false
	}

	// Check if key has expired
	if expiry, hasExpiry := db.expires[key]; hasExpiry && time.Now().After(expiry) {
		// Key has expired, remove it
		delete(db.data, key)
		delete(db.expires, key)
		return nil, false
	}

	return value, true
}

// Delete removes a key
func (db *Database) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[key]; exists {
		delete(db.data, key)
		delete(db.expires, key)
		return true
	}
	return false
}

// SetWithExpiry sets a key with an expiration time
func (db *Database) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	db.expires[key] = time.Now().Add(expiry)
}

// Keys returns all keys in the database
func (db *Database) Keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		// Check if expired
		if expiry, hasExpiry := db.expires[k]; hasExpiry && time.Now().After(expiry) {
			db.Delete(k) // Remove expired key
			continue
		}
		keys = append(keys, k)
	}
	return keys
}
