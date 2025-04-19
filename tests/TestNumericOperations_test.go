package tests

import (
	"github.com/GedisCaching/Gedis/storage"
	"testing"
)

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
