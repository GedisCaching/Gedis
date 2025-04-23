package tests

import (
	"testing"

	"github.com/GedisCaching/Gedis/gedis"
)

func TestNumericOperations(t *testing.T) {
	// Create a new Gedis instance with default config
	g, err := gedis.NewGedis(gedis.Config{})
	if err != nil {
		t.Fatalf("Failed to create Gedis instance: %v", err)
	}

	t.Run("Numeric Operations", func(t *testing.T) {
		// Test INCR
		val, err := g.Incr("counter")
		if err != nil {
			t.Errorf("INCR failed: %v", err)
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}

		// Test DECR
		val, err = g.Decr("counter")
		if err != nil {
			t.Errorf("DECR failed: %v", err)
		}
		if val != 0 {
			t.Errorf("Expected 0, got %d", val)
		}
	})
}
