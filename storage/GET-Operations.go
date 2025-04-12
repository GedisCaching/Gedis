package storage

import (
	"time"
)

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

// GETDEL Get a value and delete it in a single operation
func (db *Database) GETDEL(key string) (interface{}, bool) {
	value, exists := db.Get(key)
	if !exists {
		return nil, false
	}
	db.Delete(key)
	return value, exists
}
