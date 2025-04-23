package tests

import (
	"testing"

	"github.com/GedisCaching/Gedis/gedis"
)

func TestListOperations(t *testing.T) {
	// Create a new Gedis instance with default config
	g, err := gedis.NewGedis(gedis.Config{})
	if err != nil {
		t.Fatalf("Failed to create Gedis instance: %v", err)
	}

	// Test LPUSH
	t.Run("LPUSH Operation", func(t *testing.T) {
		length, err := g.LPush("list1", "value1", "value2", "value3")
		if err != nil {
			t.Errorf("LPUSH failed: %v", err)
		}
		if length != 3 {
			t.Errorf("Expected length 3, got %d", length)
		}
	})

	// Test RPUSH
	t.Run("RPUSH Operation", func(t *testing.T) {
		length, err := g.RPush("list2", "value1", "value2", "value3")
		if err != nil {
			t.Errorf("RPUSH failed: %v", err)
		}
		if length != 3 {
			t.Errorf("Expected length 3, got %d", length)
		}
	})

	// Test LRANGE
	t.Run("LRANGE Operation", func(t *testing.T) {
		// Reset list
		g.Delete("list1")
		g.LPush("list1", "value3", "value2", "value1")

		// Get entire list
		values, err := g.LRange("list1", 0, -1)
		if err != nil {
			t.Errorf("LRANGE failed: %v", err)
		}
		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}
		if values[0] != "value3" || values[1] != "value2" || values[2] != "value1" {
			t.Errorf("Values in wrong order, expected [value3, value2, value1], got %v", values)
		}

		// Get range
		values, err = g.LRange("list1", 0, 1)
		if err != nil {
			t.Errorf("LRANGE with range failed: %v", err)
		}
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}
		if values[0] != "value3" || values[1] != "value2" {
			t.Errorf("Range values wrong, expected [value3, value2], got %v", values)
		}
	})

	// Test LPOP
	t.Run("LPOP Operation", func(t *testing.T) {
		// Reset list
		g.Delete("list1")
		g.LPush("list1", "value3", "value2", "value1")

		val, err := g.LPop("list1")
		if err != nil {
			t.Errorf("LPOP failed: %v", err)
		}
		if val != "value3" {
			t.Errorf("Expected value3, got %v", val)
		}

		// Check remaining elements
		values, _ := g.LRange("list1", 0, -1)
		if len(values) != 2 {
			t.Errorf("Expected 2 values after LPOP, got %d", len(values))
		}
	})

	// Test RPOP
	t.Run("RPOP Operation", func(t *testing.T) {
		// Reset list
		g.Delete("list2")
		g.RPush("list2", "value1", "value2", "value3")

		val, err := g.RPop("list2")
		if err != nil {
			t.Errorf("RPOP failed: %v", err)
		}
		if val != "value3" {
			t.Errorf("Expected value3, got %v", val)
		}

		// Check remaining elements
		values, _ := g.LRange("list2", 0, -1)
		if len(values) != 2 {
			t.Errorf("Expected 2 values after RPOP, got %d", len(values))
		}

		if values[0] != "value1" || values[1] != "value2" {
			t.Errorf("Range values wrong, expected [value1, value2], got %v", values)
		}
	})
}
