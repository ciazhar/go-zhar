package main

import (
	"fmt"
	"sync"
)

type Stats struct {
	mu       sync.Mutex
	counters map[string]int
}

func (s *Stats) SnapshotBad() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counters
}

func (s *Stats) SnapshotGood() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]int, len(s.counters))
	for k, v := range s.counters {
		result[k] = v
	}
	return result
}

func main() {

	fmt.Println("Bad Example")

	stats := &Stats{
		counters: map[string]int{"a": 1, "b": 2},
	}

	// Get a snapshot of the current stats
	snapshot := stats.SnapshotBad()

	// Modify the snapshot
	snapshot["a"] = 42

	// The original stats remain unchanged
	fmt.Println(stats.counters) // Output: map[a:42 b:2]
	fmt.Println(snapshot)       // Output: map[a:42 b:2]

	fmt.Println("Good Example")

	stats2 := &Stats{
		counters: map[string]int{"a": 1, "b": 2},
	}

	// Get a snapshot of the current stats
	snapshot2 := stats2.SnapshotGood()

	// Modify the snapshot
	snapshot2["a"] = 42

	// The original stats remain unchanged
	fmt.Println(stats2.counters) // Output: map[a:1 b:2]
	fmt.Println(snapshot2)       // Output: map[a:42 b:2]
}
