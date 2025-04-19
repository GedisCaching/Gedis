package tests

import (
	"testing"

	"github.com/GedisCaching/Gedis/storage"
)

func TestSortedSetOperations(t *testing.T) {
	db := storage.NewDatabase()

	// Test ZADD
	t.Run("ZADD Operation", func(t *testing.T) {
		scoreMembers := map[string]float64{
			"member1": 1.0,
			"member2": 2.0,
			"member3": 3.0,
		}

		added := db.ZADD("myset", scoreMembers)
		if added != 3 {
			t.Errorf("Expected 3 members to be added, got %d", added)
		}

		// Test adding existing members with updated scores
		updateScoreMembers := map[string]float64{
			"member1": 10.0,
			"member4": 4.0,
		}

		added = db.ZADD("myset", updateScoreMembers)
		if added != 1 {
			t.Errorf("Expected 1 new member to be added, got %d", added)
		}
	})

	// Test ZRANGE
	t.Run("ZRANGE Operation", func(t *testing.T) {
		// Get all members without scores
		result := db.ZRANGE("myset", 0, -1, false)
		if len(result) != 4 {
			t.Errorf("Expected 4 members, got %d", len(result))
		}

		// The expected order after updates should be: member2, member3, member4, member1
		expectedOrder := []string{"member2", "member3", "member4", "member1"}
		for i, expected := range expectedOrder {
			if i < len(result) {
				if result[i] != expected {
					t.Errorf("Expected member at index %d to be %s, got %v", i, expected, result[i])
				}
			}
		}

		// Get members with scores
		resultWithScores := db.ZRANGE("myset", 0, -1, true)
		if len(resultWithScores) != 8 { // 4 members x 2 (member + score)
			t.Errorf("Expected 8 items (members and scores), got %d", len(resultWithScores))
		}

		// Test range with specific bounds
		limitedRange := db.ZRANGE("myset", 1, 2, false)
		if len(limitedRange) != 2 {
			t.Errorf("Expected 2 members, got %d", len(limitedRange))
		}
		expectedLimitedOrder := []string{"member3", "member4"}
		for i, expected := range expectedLimitedOrder {
			if limitedRange[i] != expected {
				t.Errorf("Expected member at index %d to be %s, got %v", i, expected, limitedRange[i])
			}
		}

		// Test range with non-existent key
		nonExistResult := db.ZRANGE("nonexistentkey", 0, -1, false)
		if len(nonExistResult) != 0 {
			t.Errorf("Expected empty result for non-existent key, got %d items", len(nonExistResult))
		}
	})

	// Test ZRANK
	t.Run("ZRANK Operation", func(t *testing.T) {
		// Test ranks of existing members
		rank, exists := db.ZRANK("myset", "member1")
		if !exists {
			t.Error("ZRANK failed: member1 should exist")
		}
		if rank != 3 { // should be last (index 3) after the score update to 10.0
			t.Errorf("Expected rank 3 for member1, got %d", rank)
		}

		rank, exists = db.ZRANK("myset", "member2")
		if !exists {
			t.Error("ZRANK failed: member2 should exist")
		}
		if rank != 0 { // should be first (index 0) with score 2.0
			t.Errorf("Expected rank 0 for member2, got %d", rank)
		}

		// Test rank of non-existent member
		rank, exists = db.ZRANK("myset", "nonexistentmember")
		if exists {
			t.Error("ZRANK should return false for non-existent member")
		}
		if rank != -1 {
			t.Errorf("Expected rank -1 for non-existent member, got %d", rank)
		}

		// Test rank in non-existent key
		rank, exists = db.ZRANK("nonexistentkey", "member1")
		if exists {
			t.Error("ZRANK should return false for non-existent key")
		}
		if rank != -1 {
			t.Errorf("Expected rank -1 for member in non-existent key, got %d", rank)
		}
	})
}
