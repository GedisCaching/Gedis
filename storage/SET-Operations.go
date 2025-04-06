package storage

import (
	"time"
)

// Set stores a key-value pair
func (db *Database) Set(key interface{}, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// SetWithExpiry sets a key with an expiration time
func (db *Database) SetWithExpiry(key interface{}, value interface{}, expiry time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	db.expires[key] = time.Now().Add(expiry)
}
