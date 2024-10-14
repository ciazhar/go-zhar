package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"time"
)

func startTrace() {
	f, err := os.Create("trace.out")
	if err != nil {
		fmt.Println("Failed to create trace output file:", err)
		return
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		fmt.Println("Failed to start trace:", err)
		return
	}
	defer trace.Stop()

	// Simulate workload
	for i := 0; i < 10; i++ {
		s := make([]byte, 10<<20) // Allocate 10 MiB
		if i%2 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("Iteration %d\n", i)
		_ = s // Use the allocated memory
	}
}

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Start trace
	startTrace()

	// Simulate workload
	for i := 0; i < 10; i++ {
		s := make([]byte, 10<<20) // Allocate 10 MiB
		if i%2 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Printf("Iteration %d\n", i)
		_ = s // Use the allocated memory
	}

	// Keep the program running to access pprof endpoint
	select {}
}
