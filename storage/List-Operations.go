package storage

import (
	"errors"
)

// LPush adds one or more values to the beginning of a list
func (db *Database) LPush(key string, values ...interface{}) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	var list []interface{}
	existingVal, exists := db.data[key]

	if exists {
		// If key exists, check that it's a list
		var ok bool
		list, ok = existingVal.([]interface{})
		if !ok {
			return 0, errors.New("value at key is not a list")
		}
	} else {
		// If key doesn't exist, create a new list
		list = []interface{}{}
	}

	// Prepend values to the list
	newList := make([]interface{}, len(values)+len(list))
	copy(newList[len(values):], list)
	for i, v := range values {
		newList[i] = v
	}

	// Store updated list
	db.data[key] = newList

	return len(newList), nil
}

// RPush adds one or more values to the end of a list
func (db *Database) RPush(key string, values ...interface{}) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	var list []interface{}
	existingVal, exists := db.data[key]

	if exists {
		// If key exists, check that it's a list
		var ok bool
		list, ok = existingVal.([]interface{})
		if !ok {
			return 0, errors.New("value at key is not a list")
		}
	} else {
		// If key doesn't exist, create a new list
		list = []interface{}{}
	}

	// Append values to the list
	list = append(list, values...)

	// Store updated list
	db.data[key] = list

	return len(list), nil
}

// LRange returns a range of elements from a list
func (db *Database) LRange(key string, start, stop int) ([]interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if key exists
	existingVal, exists := db.data[key]
	if !exists {
		return []interface{}{}, nil // Redis returns empty list for non-existent keys
	}

	// Check if value is a list
	list, ok := existingVal.([]interface{})
	if !ok {
		return nil, errors.New("value at key is not a list")
	}

	length := len(list)

	// Handle negative indices (counting from the end)
	if start < 0 {
		start = length + start
		if start < 0 {
			start = 0
		}
	}

	if stop < 0 {
		stop = length + stop
	}

	// Bounds checking
	if start >= length || start > stop {
		return []interface{}{}, nil
	}

	if stop >= length {
		stop = length - 1
	}

	// Return the range (inclusive of stop in Redis)
	return list[start : stop+1], nil
}
