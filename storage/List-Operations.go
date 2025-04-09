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
	copy(newList, values)

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

// Remove and return the first element of a list
func (db *Database) LPop(key string) (interface{}, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	existingVal, exists := db.data[key]
	if !exists {
		return nil, nil // Redis returns nil for non-existent keys
	}

	// Check if value is a list
	list, ok := existingVal.([]interface{})
	if !ok {
		return nil, errors.New("value at key is not a list")
	}

	if len(list) == 0 {
		return nil, nil // Return nil if the list is empty
	}

	// Remove and return the first element
	firstElement := list[0]
	list = list[1:]

	// Store updated list
	db.data[key] = list

	return firstElement, nil
}

// Remove and return the last element of a list
func (db *Database) RPop(key string) (interface{}, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	existingVal, exists := db.data[key]
	if !exists {
		return nil, nil // Redis returns nil for non-existent keys
	}

	// Check if value is a list
	list, ok := existingVal.([]interface{})
	if !ok {
		return nil, errors.New("value at key is not a list")
	}

	if len(list) == 0 {
		return nil, nil // Return nil if the list is empty
	}

	// Remove and return the last element
	lastElement := list[len(list)-1]
	list = list[:len(list)-1]

	// Store updated list
	db.data[key] = list

	return lastElement, nil
}

// Get the length of a list
func (db *Database) LLen(key string) (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if key exists
	existingVal, exists := db.data[key]
	if !exists {
		return 0, errors.New("key does not exist")
	}

	// Check if value is a list
	list, ok := existingVal.([]interface{})
	if !ok {
		return 0, errors.New("value at key is not a list")
	}

	return len(list), nil
}

// Set the value of an element in a list by its index
func (db *Database) LSet(key string, index int, value interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists
	existingVal, exists := db.data[key]
	if !exists {
		return errors.New("key does not exist")
	}

	// Check if value is a list
	list, ok := existingVal.([]interface{})
	if !ok {
		return errors.New("value at key is not a list")
	}

	// Check index bounds
	if index < 0 || index >= len(list) {
		return errors.New("index out of range")
	}

	// Set the value at the specified index
	list[index] = value

	// Store updated list
	db.data[key] = list

	return nil
}
