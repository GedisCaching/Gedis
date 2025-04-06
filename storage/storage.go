package storage

import (
	"sync"
	"time"
)

// DB represents the main database interface
type DB interface {
	// SET Operations
	Set(key interface{}, value interface{})
	SetWithExpiry(key interface{}, value interface{}, expiry time.Duration)

	// GET Operations
	Keys() []interface{}
	Get(key interface{}) (interface{}, bool)

	// DEL Operations
	Delete(key interface{}) bool

	// Incr and Decr operations
	Incr(key interface{}) (int, error)
	Decr(key interface{}) (int, error)
}

// Database represents "in-memory" Redis-like database
type Database struct {
	data    map[interface{}]interface{}
	mu      sync.RWMutex // Mutex for concurrent access
	expires map[interface{}]time.Time
}

// NewDatabase creates a new "in-memory" database
func NewDatabase() *Database {
	return &Database{
		data:    make(map[interface{}]interface{}),
		expires: make(map[interface{}]time.Time),
	}
}
