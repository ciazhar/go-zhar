package main

import (
	"fmt"
	"github.com/google/gops/agent"
	"net/http"
)

func main() {
	// Start the GOPS agent
	if err := agent.Listen(agent.Options{}); err != nil {
		fmt.Println("Error starting GOPS agent:", err)
		return
	}

	// Set up a simple HTTP server to simulate some work
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Run the HTTP server in a separate goroutine
	//go func() {
	fmt.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
	//}()

	//// Simulate the main process doing some work
	//for i := 0; i < 10; i++ {
	//	fmt.Println("Working...", i)
	//	time.Sleep(1 * time.Second)
	//}

	//fmt.Println("Program finished")
}
