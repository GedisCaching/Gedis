package storage

import (
	"time"
)

// TTL returns the remaining time to live for a key
// Returns:
// - Positive duration: remaining time until expiration
// - Zero duration: the key exists but has no expiry
// - bool=false: the key doesn't exist
func (db *Database) TTL(key string) (time.Duration, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if key exists
	_, exists := db.data[key]

	if !exists {
		return 0, false
	}

	// Check if key has an expiry
	expiry, hasExpiry := db.expires[key]
	if !hasExpiry {
		return 0, true // Key exists but has no expiry
	}

	// Calculate remaining time
	now := time.Now()

	// If the key is already expired, it should be deleted
	// (this would normally happen on access, but we'll check here too)
	if now.After(expiry) {
		// We need to upgrade to a write lock to delete
		db.mu.RUnlock()
		db.mu.Lock()
		defer db.mu.Unlock()

		// Check again after acquiring write lock
		if e, exists := db.expires[key]; exists && now.After(e) {
			delete(db.data, key)
			delete(db.expires, key)
		}
		return 0, false
	}

	// Return time until expiry
	return expiry.Sub(now), true
}
