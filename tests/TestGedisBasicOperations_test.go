package tests

import (
	"testing"

	"github.com/GedisCaching/Gedis/gedis"
)

func TestBasicOperations(t *testing.T) {
	// Create a new Gedis instance with default config
	g, err := gedis.NewGedis(gedis.Config{})
	if err != nil {
		t.Fatalf("Failed to create Gedis instance: %v", err)
	}

	// Test SET
	t.Run("SET Operation", func(t *testing.T) {
		g.Set("key1", "value1")
	})

	// Test GET
	t.Run("GET Operation", func(t *testing.T) {
		val, exists := g.Get("key1")
		if !exists {
			t.Errorf("GET failed: key not found")
		}
		if val != "value1" {
			t.Errorf("Expected value1, got %v", val)
		}
	})

	// Test DEL
	t.Run("DEL Operation", func(t *testing.T) {
		success := g.Delete("key1")
		if !success {
			t.Errorf("DEL failed")
		}
		_, exists := g.Get("key1")
		if exists {
			t.Error("Key should be deleted")
		}
	})

	// Test GETDEL
	t.Run("GETDEL Operation", func(t *testing.T) {
		g.Set("key2", "value2")
		val, exists := g.GETDEL("key2")
		if !exists {
			t.Error("GETDEL failed: key not found")
		}
		if val != "value2" {
			t.Errorf("Expected value2, got %v", val)
		}
		_, exists = g.Get("key2")
		if exists {
			t.Error("Key should be deleted after GETDEL")
		}
	})

	// Test RENAME
	t.Run("RENAME Operation", func(t *testing.T) {
		g.Set("original", "value")
		err := g.RENAME("original", "renamed")
		if err != nil {
			t.Errorf("RENAME failed: %v", err)
		}
		_, exists := g.Get("original")
		if exists {
			t.Error("Original key should not exist after rename")
		}
		val, exists := g.Get("renamed")
		if !exists {
			t.Error("New key should exist after rename")
		}
		if val != "value" {
			t.Errorf("Expected value, got %v", val)
		}
	})

	// Test KEYS
	t.Run("KEYS Operation", func(t *testing.T) {
		g.Set("key1", "value1")
		g.Set("key2", "value2")
		g.Set("key3", "value3")

		keys := g.Keys()
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
