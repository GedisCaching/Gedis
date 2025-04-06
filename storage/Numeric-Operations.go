package storage

import (
	"fmt"
	"strconv"
)

// Incr increments the value of a key by one.
// If the key has expired, it is deleted from the database
// and a new key is created with the incremented value and no expiration time.
func (db *Database) Incr(key string) (int, error) {
	value, exists := db.Get(key)

	// Acquire the lock after Get because the Get operation already applies its own lock,
	// and we cannot acquire a lock that is already held.
	db.mu.Lock()
	defer db.mu.Unlock()

	var intValue int
	if !exists {
		// In Redis, INCR creates the key with value 0 if it doesn't exist, then increments it
		intValue = 1
	} else {
		switch v := value.(type) {
		case int:
			intValue = v + 1
		case string:
			if parsedInt, err := strconv.Atoi(v); err == nil {
				intValue = parsedInt + 1
			} else {
				return 0, fmt.Errorf("invalid value for key %s: %s", key, v)
			}
		default:
			return 0, fmt.Errorf("invalid value type for key %s: %T", key, value)
		}
	}

	db.data[key] = intValue
	return intValue, nil
}

// Decr decrements the value of a key by one.
// If the key has expired, it is deleted from the database
// and a new key is created with the decremented value and no expiration time.
func (db *Database) Decr(key string) (int, error) {
	value, exists := db.Get(key)

	// Acquire the lock after Get because the Get operation already applies its own lock,
	// and we cannot acquire a lock that is already held.
	db.mu.Lock()
	defer db.mu.Unlock()

	var intValue int
	if !exists {
		// In Redis, DECR creates the key with value 0 if it doesn't exist, then decrements it
		intValue = -1
	} else {
		switch v := value.(type) {
		case int:
			intValue = v - 1
		case string:
			if parsedInt, err := strconv.Atoi(v); err == nil {
				intValue = parsedInt - 1
			} else {
				return 0, fmt.Errorf("invalid value for key %s: %s", key, v)
			}
		default:
			return 0, fmt.Errorf("invalid value type for key %s: %T", key, value)
		}
	}

	db.data[key] = intValue
	return intValue, nil
}
