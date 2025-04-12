package storage

import (
	"errors"
	"time"
)

// Set stores a key-value pair
func (db *Database) Set(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

// SetWithExpiry sets a key with an expiration time
func (db *Database) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	db.expires[key] = time.Now().Add(expiry)
}

// DEXPIRE set expiration on existing key
func (db *Database) DEXPIRE(key string, expiry time.Duration) error {
	_, exists := db.Get(key)
	if !exists {
		return errors.New("key does not exist")
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.expires[key] = time.Now().Add(expiry)
	return nil
}
