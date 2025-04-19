package tests

import (
	"github.com/GedisCaching/Gedis/storage"
	"testing"
)

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
