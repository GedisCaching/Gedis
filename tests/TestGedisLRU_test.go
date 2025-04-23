package tests

import (
	"testing"
	"time"

	"github.com/GedisCaching/Gedis/gedis"
	redis "github.com/GedisCaching/Gedis/server"
)

func TestLRUFunctionality(t *testing.T) {
	// Reset to default capacity of 100 to clear any previous state
	redis.SetGlobalCapacity(100)

	// Save the initial count to track our progress
	initialCount := redis.GetGlobalServerCount()
	t.Logf("Initial server count: %d", initialCount)

	// Set a small capacity for testing
	redis.SetGlobalCapacity(2)

	// Create servers with unique addresses for testing
	addresses := []string{
		"localhost:8001",
		"localhost:8002",
		"localhost:8003",
		"localhost:8004",
	}

	// Create and store the first two (which should fill capacity)
	g1, err1 := gedis.NewGedis(gedis.Config{Address: addresses[0]})
	if err1 != nil {
		t.Fatalf("Failed to create first Gedis instance: %v", err1)
	}

	g2, err2 := gedis.NewGedis(gedis.Config{Address: addresses[1]})
	if err2 != nil {
		t.Fatalf("Failed to create second Gedis instance: %v", err2)
	}

	// Use each instance to ensure they're cached
	g1.Set("key1", "value1")
	g2.Set("key2", "value2")

	// Check server count
	afterTwoCount := redis.GetGlobalServerCount()
	t.Logf("After creating 2 servers: %d servers", afterTwoCount)

	// Create a third instance - this should cause eviction of the least recently used (g1)
	g3, err3 := gedis.NewGedis(gedis.Config{Address: addresses[2]})
	if err3 != nil {
		t.Fatalf("Failed to create third Gedis instance: %v", err3)
	}
	g3.Set("key3", "value3")

	// Check server count again - should still be 2 (one in, one out)
	afterThreeCount := redis.GetGlobalServerCount()
	t.Logf("After creating 3rd server: %d servers", afterThreeCount)

	if afterThreeCount > 2 {
		t.Errorf("LRU eviction failed. Expected 2 servers, got %d", afterThreeCount)
	}

	// Use g2 again to ensure it stays as the most recently used
	g2.Get("key2")
	time.Sleep(10 * time.Millisecond)

	// Create a fourth instance - this should evict g3 (since g2 was just used)
	g4, err4 := gedis.NewGedis(gedis.Config{Address: addresses[3]})
	if err4 != nil {
		t.Fatalf("Failed to create fourth Gedis instance: %v", err4)
	}
	g4.Set("key4", "value4")

	// Check count again
	afterFourCount := redis.GetGlobalServerCount()
	t.Logf("After creating 4th server: %d servers", afterFourCount)

	if afterFourCount > 2 {
		t.Errorf("LRU eviction failed. Expected 2 servers, got %d", afterFourCount)
	}

	// Restore default capacity
	redis.SetGlobalCapacity(100)
}
