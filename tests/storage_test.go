package tests

import (
	"github.com/GedisCaching/Gedis/storage"
	"testing"
	"time"
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
