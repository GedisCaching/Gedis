package storage

import (
	"fmt"
)

// ------------------------------ Hash Operations ------------------------------

// HSET sets the value of a field in a hash
func (db *Database) HSET(key string, field string, value interface{}) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Initialize the hash if it doesn't exist
	if _, exists := db.data[key]; !exists {
		db.data[key] = make(map[string]interface{})
	}

	hash, ok := db.data[key].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("key %s is not a hash", key)
	}

	hash[field] = value
	return true, nil
}

// HGET retrieves the value of a field in a hash
func (db *Database) HGET(key string, field string) (interface{}, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return nil, false
	}

	value, exists := hash[field]
	return value, exists
}

// HDEL deletes a field from a hash
func (db *Database) HDEL(key string, field string) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return false, fmt.Errorf("key %s is not a hash", key)
	}

	if _, exists := hash[field]; !exists {
		return false, nil
	}

	delete(hash, field)
	return true, nil
}

// HGETALL retrieves all fields and values in a hash
func (db *Database) HGETALL(key string) (map[string]interface{}, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return nil, false
	}

	// Create a copy of the hash to avoid concurrent modification
	result := make(map[string]interface{})
	for k, v := range hash {
		result[k] = v
	}
	return result, true
}

// HKEYS retrieves all field names in a hash
func (db *Database) HKEYS(key string) ([]string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return nil, false
	}

	keys := make([]string, 0, len(hash))
	for k := range hash {
		keys = append(keys, k)
	}
	return keys, true
}

// HVALS retrieves all values in a hash
func (db *Database) HVALS(key string) ([]interface{}, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return nil, false
	}

	values := make([]interface{}, 0, len(hash))
	for _, v := range hash {
		values = append(values, v)
	}
	return values, true
}

// HLEN retrieves the number of fields in a hash
func (db *Database) HLEN(key string) (int, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	hash, exists := db.data[key].(map[string]interface{})
	if !exists {
		return 0, false
	}

	return len(hash), true
}
