package tests

import (
	"testing"
	"time"

	"github.com/GedisCaching/Gedis/storage"
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

func TestTTLOperations(t *testing.T) {
	db := storage.NewDatabase()

	t.Run("TTL Operations", func(t *testing.T) {
		// Set with TTL
		db.SetWithExpiry("key1", "value1", 2*time.Second)

		// Check TTL
		ttl, exist := db.TTL("key1")
		if !exist {
			t.Error("key1 not exist")
		}

		if ttl < 0 {
			t.Error("TTL should be positive")
		}

		// Wait for expiration
		time.Sleep(3 * time.Second)

		// Check if key is expired
		_, exists := db.Get("key1")
		if exists == true {
			t.Error("Key should be expired")
		}
	})

	// Test DEXPIRE
	t.Run("DEXPIRE Operation", func(t *testing.T) {
		db.Set("key2", "value2")

		err := db.DEXPIRE("key2", 1*time.Second)
		if err != nil {
			t.Errorf("DEXPIRE failed: %v", err)
		}

		// Verify key still exists
		_, exists := db.Get("key2")
		if !exists {
			t.Error("Key should exist immediately after DEXPIRE")
		}

		// Wait for expiration
		time.Sleep(2 * time.Second)

		// Check if key is expired
		_, exists = db.Get("key2")
		if exists {
			t.Error("Key should be expired after DEXPIRE timeout")
		}
	})
}

func TestListOperations(t *testing.T) {
	db := storage.NewDatabase()

	t.Run("List Operations", func(t *testing.T) {
		// Test LPUSH
		length, err := db.LPush("list1", "value1", "value2", "value3")
		if err != nil {
			t.Errorf("LPUSH failed: %v", err)
		}
		if length != 3 {
			t.Errorf("Expected length 3, got %d", length)
		}

		// Test LRANGE
		values, err := db.LRange("list1", 0, -1)
		if err != nil {
			t.Errorf("LRANGE failed: %v", err)
		}
		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}

		// Test RPUSH
		length, err = db.RPush("list1", "value4", "value5")
		if err != nil {
			t.Errorf("RPUSH failed: %v", err)
		}
		if length != 5 {
			t.Errorf("Expected length 5, got %d", length)
		}

		// Test updated LRANGE
		values, err = db.LRange("list1", 0, -1)
		if err != nil {
			t.Errorf("LRANGE failed: %v", err)
		}
		if len(values) != 5 {
			t.Errorf("Expected 5 values, got %d", len(values))
		}

		// Test LLEN
		length, err = db.LLen("list1")
		if err != nil {
			t.Errorf("LLEN failed: %v", err)
		}
		if length != 5 {
			t.Errorf("Expected length 5, got %d", length)
		}

		// Test LPOP
		val, err := db.LPop("list1")
		if err != nil {
			t.Errorf("LPOP failed: %v", err)
		}
		if val != values[0] {
			t.Errorf("LPOP expected %v, got %v", values[0], val)
		}

		// Verify length after LPOP
		length, err = db.LLen("list1")
		if err != nil {
			t.Errorf("LLEN failed: %v", err)
		}
		if length != 4 {
			t.Errorf("Expected length 4 after LPOP, got %d", length)
		}

		// Test RPOP
		_, err = db.RPop("list1")
		if err != nil {
			t.Errorf("RPOP failed: %v", err)
		}

		// Verify length after RPOP
		length, err = db.LLen("list1")
		if err != nil {
			t.Errorf("LLEN failed: %v", err)
		}
		if length != 3 {
			t.Errorf("Expected length 3 after RPOP, got %d", length)
		}

		// Test LSET
		err = db.LSet("list1", 0, "updated-value")
		if err != nil {
			t.Errorf("LSET failed: %v", err)
		}

		// Verify LSET worked by checking the first element
		values, err = db.LRange("list1", 0, 0)
		if err != nil {
			t.Errorf("LRANGE after LSET failed: %v", err)
		}
		if len(values) != 1 || values[0] != "updated-value" {
			t.Errorf("LSET failed: Expected updated-value, got %v", values[0])
		}
	})
}

func TestNumericOperations(t *testing.T) {
	db := storage.NewDatabase()

	t.Run("Numeric Operations", func(t *testing.T) {
		// Test INCR
		val, err := db.Incr("counter")
		if err != nil {
			t.Errorf("INCR failed: %v", err)
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}

		// Test DECR
		val, err = db.Decr("counter")
		if err != nil {
			t.Errorf("DECR failed: %v", err)
		}
		if val != 0 {
			t.Errorf("Expected 0, got %d", val)
		}
	})
}
