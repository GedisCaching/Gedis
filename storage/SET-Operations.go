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

// RENAME renames the existing key to a new key
func (db *Database) RENAME(KeyOld, KeyNew string) error {
	// Get the value of the old key
	value, exists := db.Get(KeyOld)
	if !exists {
		return errors.New("key does not exist")
	}

	// Get the TTL of the old key
	ttl, exists := db.TTL(KeyOld)
	if !exists {
		return errors.New("key does not exist")
	}

	// Set the new key with the same value and TTL
	db.SetWithExpiry(KeyNew, value, ttl)

	// Delete the old key
	db.Delete(KeyOld)
	return nil
}
