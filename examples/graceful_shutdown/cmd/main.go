package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new context for handling shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize Fiber app
	app := fiber.New()

	// Define a simple route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Fiber server on :8080...")
		if err := app.Listen(":8080"); err != nil {
			fmt.Printf("Fiber server failed: %v\n", err)
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

	// Start graceful shutdown process with a 10-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown Fiber app
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down Fiber server: %v\n", err)
	} else {
		fmt.Println("Fiber server shut down successfully.")
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
