package storage

import (
	"sync"
	"time"
)

// DB represents the main database interface
type DB interface {
	// SET Operations
	Set(key string, value interface{})
	SetWithExpiry(key string, value interface{}, expiry time.Duration)

	// GET Operations
	Keys() []string
	Get(key string) (interface{}, bool)

	// DEL Operations
	Delete(key string) bool

	// List Operations
	LPush(key string, values ...interface{}) (int, error)
	RPush(key string, values ...interface{}) (int, error)
	LRange(key string, start, stop int) ([]interface{}, error)

	// TTL Operation
	TTL(key string) (time.Duration, bool)
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
