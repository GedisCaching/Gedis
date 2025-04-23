package tests

import (
	"reflect"
	"testing"

	"github.com/GedisCaching/Gedis/gedis"
)

func TestSortedSetOperations(t *testing.T) {
	// Create a new Gedis instance with default config
	g, err := gedis.NewGedis(gedis.Config{})
	if err != nil {
		t.Fatalf("Failed to create Gedis instance: %v", err)
	}

	// Test ZADD
	t.Run("ZADD Operation", func(t *testing.T) {
		// Add multiple members at once
		scoreMembers := map[string]float64{
			"member1": 1.0,
			"member2": 2.0,
			"member3": 3.0,
		}

		added := g.ZAdd("zset1", scoreMembers)
		if added != 3 {
			t.Errorf("Expected 3 members added, got %d", added)
		}

		// Add/update some more
		scoreMembers = map[string]float64{
			"member2": 2.5, // Update score
			"member4": 4.0, // New member
		}

		added = g.ZAdd("zset1", scoreMembers)
		if added != 1 {
			t.Errorf("Expected 1 new member added, got %d", added)
		}
	})

	// Test ZRANGE
	t.Run("ZRANGE Operation", func(t *testing.T) {
		// Get all members without scores
		result := g.ZRange("zset1", 0, -1, false)
		if len(result) != 4 {
			t.Errorf("Expected 4 members, got %d", len(result))
		}

		// Adjust expected order to match actual implementation
		expectedOrder := []interface{}{"member1", "member2", "member3", "member4"}
		if !reflect.DeepEqual(result, expectedOrder) {
			t.Errorf("Expected %v, got %v", expectedOrder, result)
		}

		// Get all members with scores
		resultWithScores := g.ZRange("zset1", 0, -1, true)
		if len(resultWithScores) != 8 { // 4 members * 2 (member, score)
			t.Errorf("Expected 8 items (members and scores), got %d", len(resultWithScores))
		}

		// Adjust expected scores to match actual implementation
		expectedWithScores := []interface{}{
			"member1", float64(1.0),
			"member2", float64(2.5),
			"member3", float64(3.0),
			"member4", float64(4.0),
		}

		if !reflect.DeepEqual(resultWithScores, expectedWithScores) {
			t.Errorf("Expected %v, got %v", expectedWithScores, resultWithScores)
		}

		// Adjust expected subset to match actual implementation
		subset := g.ZRange("zset1", 1, 2, false)
		expectedSubset := []interface{}{"member2", "member3"}
		if !reflect.DeepEqual(subset, expectedSubset) {
			t.Errorf("Expected subset %v, got %v", expectedSubset, subset)
		}
	})

	// Test ZRANK
	t.Run("ZRANK Operation", func(t *testing.T) {
		// Get rank of member1 (should be 0, lowest score)
		rank, exists := g.ZRank("zset1", "member1")
		if !exists {
			t.Error("ZRANK failed: member not found")
		}
		if rank != 0 {
			t.Errorf("Expected rank 0 for member1, got %d", rank)
		}

		// Get rank of member4 (should be 3, highest score)
		rank, exists = g.ZRank("zset1", "member4")
		if !exists {
			t.Error("ZRANK failed: member not found")
		}
		if rank != 3 {
			t.Errorf("Expected rank 3 for member4, got %d", rank)
		}

		// Get rank of non-existent member
		rank, exists = g.ZRank("zset1", "nonexistent")
		if exists {
			t.Error("ZRANK should return false for non-existent member")
		}

		if rank != -1 {
			t.Errorf("Expected rank -1 for non-existent member, got %d", rank)
		}
	})
}
