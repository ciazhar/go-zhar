package main

import (
	"fmt"
	"runtime"
	"time"
)

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func main() {
	// Print initial memory stats
	printMemStats()

	// Simulate workload
	for i := 0; i < 10; i++ {
		s := make([]byte, 10<<20) // Allocate 10 MiB
		if i%2 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("Iteration %d\n", i)
		printMemStats()
		_ = s // Use the allocated memory
	}

	// Final memory stats
	printMemStats()
}
