package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Create a new context for handling shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the HTTP server in a separate goroutine
	server := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting server...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Server failed: %v\n", err)
		}
	}()

	// Simulate a background worker with a long-running task
	wg.Add(1)
	go func() {
		defer wg.Done()
		backgroundWorker(ctx)
	}()

	// Listen for shutdown signal
	<-ctx.Done() // Block until signal is received
	fmt.Println("\nShutting down gracefully...")

	// Start graceful shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	} else {
		fmt.Println("Server shut down successfully.")
	}

	// Wait for background tasks to complete
	wg.Wait()
	fmt.Println("All goroutines completed.")
}

func backgroundWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // Listen for shutdown signal
			fmt.Println("Background worker shutting down...")
			return
		default:
			// Simulate work
			fmt.Println("Background worker is working...")
			time.Sleep(2 * time.Second)
		}
	}
}
