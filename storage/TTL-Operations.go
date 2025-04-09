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
	// First, check existence and expiry using a read lock.
	db.mu.RLock()
	_, exists := db.data[key]
	expiry, hasExpiry := db.expires[key]
	db.mu.RUnlock()

	// If the key doesn't exist, return false.
	if !exists {
		return 0, false
	}

	// If there's no expiry, return zero duration.
	if !hasExpiry {
		return 0, true
	}

	now := time.Now()
	// If expired, acquire a write lock to delete the key.
	if now.After(expiry) {
		db.mu.Lock()
		// Double-check the expiry and existence now.
		if exp, exists := db.expires[key]; exists && now.After(exp) {
			delete(db.data, key)
			delete(db.expires, key)
			db.mu.Unlock()
			return 0, false
		}
		db.mu.Unlock()
	}

	// If not expired, return the remaining time.
	return expiry.Sub(now), true
}
