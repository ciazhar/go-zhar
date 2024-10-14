package main

import (
	"fmt"
	"sync"
)

type Trip struct {
	destination string
}

type Driver struct {
	trips []Trip
}

//func (d *Driver) SetTrips(trips []Trip) {
//	d.trips = trips
//}

func (d *Driver) SetTrips(trips []Trip) {
	// Create a copy of the trips slice
	d.trips = make([]Trip, len(trips))
	copy(d.trips, trips)
}

type Stats struct {
	mu       sync.Mutex
	counters map[string]int
}

// Snapshot returns the current stats.
func (s *Stats) Snapshot() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counters
}

func main() {
	stats := &Stats{
		counters: map[string]int{"a": 1, "b": 2},
	}

	// Get a snapshot of the current stats
	snapshot := stats.Snapshot()

	// Modify the snapshot
	snapshot["a"] = 42

	// The original stats remain unchanged
	fmt.Println(stats.counters) // Output: map[a:1 b:2]
	fmt.Println(snapshot)       // Output: map[a:42 b:2]
}
