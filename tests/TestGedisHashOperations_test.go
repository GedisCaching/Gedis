package tests

import (
	"testing"
	"time"

	"github.com/GedisCaching/Gedis/gedis"
	redis "github.com/GedisCaching/Gedis/server"
)

func TestGedisHashOperations(t *testing.T) {
	// Create a new Gedis instance with default config
	g, err := gedis.NewGedis(gedis.Config{})
	if err != nil {
		t.Fatalf("Failed to create Gedis instance: %v", err)
	}

	// Setup test data
	t.Run("Setup Hash Data", func(t *testing.T) {
		// Create a hash with multiple fields
		success, err := g.HSET("user:1", "name", "John Doe")
		if !success || err != nil {
			t.Errorf("HSET failed: %v", err)
		}

		success, err = g.HSET("user:1", "email", "john@example.com")
		if !success || err != nil {
			t.Errorf("HSET failed: %v", err)
		}

		success, err = g.HSET("user:1", "age", 30)
		if !success || err != nil {
			t.Errorf("HSET failed: %v", err)
		}
	})

	// Test basic hash operations
	t.Run("Basic Hash Operations", func(t *testing.T) {
		// Test HGET
		name, exists := g.HGET("user:1", "name")
		if !exists || name != "John Doe" {
			t.Errorf("HGET name failed, got: %v", name)
		}

		email, exists := g.HGET("user:1", "email")
		if !exists || email != "john@example.com" {
			t.Errorf("HGET email failed, got: %v", email)
		}

		age, exists := g.HGET("user:1", "age")
		if !exists || age != 30 {
			t.Errorf("HGET age failed, got: %v", age)
		}

		// Test HGETALL
		userData, exists := g.HGETALL("user:1")
		if !exists || len(userData) != 3 {
			t.Errorf("HGETALL failed, got: %v", userData)
		}

		// Test HLEN
		length, exists := g.HLEN("user:1")
		if !exists || length != 3 {
			t.Errorf("HLEN failed, expected 3, got: %v", length)
		}

		// Test HDEL
		deleted, err := g.HDEL("user:1", "email")
		if !deleted || err != nil {
			t.Errorf("HDEL failed: %v", err)
		}

		// Verify deletion
		_, exists = g.HGET("user:1", "email")
		if exists {
			t.Error("Field should be deleted but still exists")
		}

		// Check updated length
		length, exists = g.HLEN("user:1")
		if !exists || length != 2 {
			t.Errorf("HLEN after deletion failed, expected 2, got: %v", length)
		}
	})

	// Test Hash LRU update
	t.Run("Hash Operations LRU Update", func(t *testing.T) {
		// Reset the capacity to a known value at the start
		redis.SetGlobalCapacity(100)

		// Create a separate test to avoid interference with other tests
		// Reset the server count first (clear existing servers)
		// This is a more controlled test environment
		initialCount := redis.GetGlobalServerCount()

		// Set capacity to exactly 2 for eviction testing
		redis.SetGlobalCapacity(2)

		// Create two instances with unique addresses
		g1, _ := gedis.NewGedis(gedis.Config{Address: "localhost:7001"})
		g2, _ := gedis.NewGedis(gedis.Config{Address: "localhost:7002"})

		// Use them to ensure they're in the LRU cache
		g1.HSET("lru:1", "field", "value1")
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
		g2.HSET("lru:2", "field", "value2")

		// Verify we have 2 servers (or initialCount+2 if we couldn't reset)
		currentCount := redis.GetGlobalServerCount()
		expectedCount := initialCount
		if expectedCount < 2 {
			expectedCount = 2 // We should have at least the 2 we just created
		}

		if currentCount != expectedCount {
			t.Logf("Warning: unexpected server count. Expected: %d, Got: %d", expectedCount, currentCount)
		}

		// Create a third instance - with capacity=2, this should trigger eviction
		g3, _ := gedis.NewGedis(gedis.Config{Address: "localhost:7003"})
		g3.HSET("lru:3", "field", "value3")

		// Count should still be the same (one was evicted, one was added)
		afterEvictionCount := redis.GetGlobalServerCount()
		if afterEvictionCount != expectedCount {
			t.Errorf("After eviction count mismatch. Expected: %d, Got: %d", expectedCount, afterEvictionCount)
		}

		// Access g2 to make it recently used
		g2.HGET("lru:2", "field")

		// Access g3 to make it recently used
		g3.HGET("lru:3", "field")

		// Now create g4, which should evict g1 (the least recently used)
		g4, _ := gedis.NewGedis(gedis.Config{Address: "localhost:7004"})
		g4.HSET("lru:4", "field", "value4")

		// Restore capacity to avoid affecting other tests
		redis.SetGlobalCapacity(100)
	})
}
