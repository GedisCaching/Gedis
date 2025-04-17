package storage

import (
	"sort"
)

// SortedSetItem represents a member and score in a sorted set
type SortedSetItem struct {
	Member string
	Score  float64
}

// SortedSet represents a sorted set data structure
type SortedSet struct {
	Items []SortedSetItem
}

// NewSortedSet creates a new sorted set
func NewSortedSet() *SortedSet {
	return &SortedSet{
		Items: []SortedSetItem{},
	}
}

// Add adds or updates a member with a score in the sorted set
// Returns true if the member was added (didn't exist before)
func (ss *SortedSet) Add(score float64, member string) bool {
	// Check if member already exists
	for i := range ss.Items {
		if ss.Items[i].Member == member {
			ss.Items[i].Score = score
			// Re-sort the set
			sort.Slice(ss.Items, func(i, j int) bool {
				return ss.Items[i].Score < ss.Items[j].Score
			})
			return false
		}
	}

	// Add new member
	ss.Items = append(ss.Items, SortedSetItem{
		Member: member,
		Score:  score,
	})

	// Sort the set by score
	sort.Slice(ss.Items, func(i, j int) bool {
		return ss.Items[i].Score < ss.Items[j].Score
	})

	return true
}

// Range returns members in the range [start, stop]
// If withScores is true, returns members and scores alternating
func (ss *SortedSet) Range(start, stop int, withScores bool) []interface{} {
	// Handle negative indices (counting from the end)
	if start < 0 {
		start = len(ss.Items) + start
	}
	if stop < 0 {
		stop = len(ss.Items) + stop
	}

	// Ensure indices are within bounds
	if start < 0 {
		start = 0
	}
	if stop >= len(ss.Items) {
		stop = len(ss.Items) - 1
	}
	if start > stop || start >= len(ss.Items) {
		return []interface{}{}
	}

	// Build result array
	result := make([]interface{}, 0)
	for i := start; i <= stop; i++ {
		result = append(result, ss.Items[i].Member)
		if withScores {
			result = append(result, ss.Items[i].Score)
		}
	}

	return result
}

// Rank returns the rank of a member (0-based) or -1 if not found
func (ss *SortedSet) Rank(member string) int {
	for i, item := range ss.Items {
		if item.Member == member {
			return i
		}
	}
	return -1
}

// ZADD adds one or more members to a sorted set, or updates their score if already exist
func (db *Database) ZADD(key string, scoreMembers map[string]float64) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if key exists and is a sorted set
	val, exists := db.setStorage[key]
	if !exists {
		val = NewSortedSet()
		db.setStorage[key] = val
	}

	// Add members to the sorted set
	count := 0
	for score, member := range scoreMembers {
		if val.Add(member, score) {
			count++
		}
	}

	// Update the sorted set in the storage
	db.setStorage[key] = val
	return count, nil
}

// ZRANGE returns a range of members in a sorted set, by index
func (db *Database) ZRANGE(key string, start, stop int, withScores bool) ([]interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if key exists and is a sorted set
	val, exists := db.setStorage[key]
	if !exists {
		return []interface{}{}, nil
	}

	return val.Range(start, stop, withScores), nil
}

// ZRANK returns the rank of a member in a sorted set, with scores ordered from low to high
func (db *Database) ZRANK(key, member string) (int, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if key exists and is a sorted set
	val, exists := db.setStorage[key]
	if !exists {
		return -1, false
	}

	// Get the rank of the member
	rank := val.Rank(member)
	if rank == -1 {
		return -1, false
	}

	return rank, true
}
