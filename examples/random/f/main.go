package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	// Set GOGC to adjust the garbage collection aggressiveness.
	// For example, setting it to 200 will make the GC run less frequently.
	runtime.GOMAXPROCS(1)
	fmt.Println("Setting GOGC to 200")
	_ = os.Setenv("GOGC", "200")

	// Run a function that allocates memory to trigger the GC.
	allocateMemory()

	// Monitor GC performance using runtime.ReadMemStats.
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("Alloc: %v MiB\n", memStats.Alloc/1024/1024)
	fmt.Printf("TotalAlloc: %v MiB\n", memStats.TotalAlloc/1024/1024)
	fmt.Printf("Sys: %v MiB\n", memStats.Sys/1024/1024)
	fmt.Printf("NumGC: %v\n", memStats.NumGC)

	// You can monitor these stats over time in your application to see how changes in GOGC affect performance.
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		runtime.ReadMemStats(&memStats)
		fmt.Printf("Time: %v, NumGC: %v, Alloc: %v MiB\n", i, memStats.NumGC, memStats.Alloc/1024/1024)
	}
}

func allocateMemory() {
	var memoryHog []byte
	for i := 0; i < 10; i++ {
		memoryHog = append(memoryHog, make([]byte, 10*1024*1024)...)
		time.Sleep(500 * time.Millisecond)
	}
}
