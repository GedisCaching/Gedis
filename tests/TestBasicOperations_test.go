package tests

import (
	"github.com/GedisCaching/Gedis/storage"
	"testing"
)

func TestBasicOperations(t *testing.T) {
	db := storage.NewDatabase()

	// Test SET
	t.Run("SET Operation", func(t *testing.T) {
		db.Set("key1", "value1")
	})

	// Test GET
	t.Run("GET Operation", func(t *testing.T) {
		val, err := db.Get("key1")
		if err != true {
			t.Errorf("GET failed: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected value1, got %v", val)
		}
	})

	// Test DEL
	t.Run("DEL Operation", func(t *testing.T) {
		err := db.Delete("key1")
		if err != true {
			t.Errorf("DEL failed: %v", err)
		}
		_, err = db.Get("key1")
		if err == true {
			t.Error("Key should be deleted")
		}
	})

	// Test GETDEL
	t.Run("GETDEL Operation", func(t *testing.T) {
		db.Set("key2", "value2")
		val, exists := db.GETDEL("key2")
		if !exists {
			t.Error("GETDEL failed: key not found")
		}
		if val != "value2" {
			t.Errorf("Expected value2, got %v", val)
		}
		_, exists = db.Get("key2")
		if exists {
			t.Error("Key should be deleted after GETDEL")
		}
	})

	// Test RENAME
	t.Run("RENAME Operation", func(t *testing.T) {
		db.Set("original", "value")
		err := db.RENAME("original", "renamed")
		if err != nil {
			t.Errorf("RENAME failed: %v", err)
		}
		_, exists := db.Get("original")
		if exists {
			t.Error("Original key should not exist after rename")
		}
		val, exists := db.Get("renamed")
		if !exists {
			t.Error("New key should exist after rename")
		}
		if val != "value" {
			t.Errorf("Expected value, got %v", val)
		}
	})

	// Test KEYS
	t.Run("KEYS Operation", func(t *testing.T) {
		db.Set("key1", "value1")
		db.Set("key2", "value2")
		db.Set("key3", "value3")

		keys := db.Keys()
		if len(keys) < 3 {
			t.Errorf("Expected at least 3 keys, got %d", len(keys))
		}

		// Check if all keys exist in the result
		keyMap := make(map[string]bool)
		for _, key := range keys {
			keyMap[key] = true
		}

		if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
			t.Error("Not all keys were returned by KEYS operation")
		}
	})
}
