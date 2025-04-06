package storage

// Delete removes a key
func (db *Database) Delete(key interface{}) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.data[key]; exists {
		delete(db.data, key)
		delete(db.expires, key)
		return true
	}
	return false
}
