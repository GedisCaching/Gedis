package tests

import (
	"github.com/GedisCaching/Gedis/storage"
	"testing"
	"time"
)

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
