package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/pprof"
	"time"
)

type BigStruct struct {
	data [1 << 20]int
}

func allocateBigStructs(n int) {
	for i := 0; i < n; i++ {
		_ = new(BigStruct)
	}
}

func main() {
	// Start an HTTP server for pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Allocate some BigStructs
	allocateBigStructs(10)

	// Allow some time for profiling
	time.Sleep(30 * time.Second)

	// Print memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc: %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("Sys: %v MB\n", m.Sys/1024/1024)
	fmt.Printf("NumGC: %v\n", m.NumGC)

	// Create a memory profile
	f, err := pprof.Create("memprofile")
	if err != nil {
		fmt.Println("Could not create memory profile: ", err)
		return
	}
	defer f.Close()
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Println("Could not write memory profile: ", err)
		return
	}
}
