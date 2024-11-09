package main

import (
	"context"
	"github.com/ciazhar/go-start-small/pkg/logger"
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
		logger.LogInfo(ctx, "Starting Fiber server on :8080...", nil)
		if err := app.Listen(":8080"); err != nil {
			logger.LogFatal(ctx, err, "Failed to start Fiber server", nil)
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
	logger.LogInfo(ctx, "Received shutdown signal...", nil)

	// Start graceful shutdown process with a 10-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown Fiber app
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.LogFatal(ctx, err, "Error shutting down Fiber server", nil)
	} else {
		logger.LogInfo(ctx, "Fiber server shut down gracefully.", nil)
	}

	// Wait for background tasks to complete
	wg.Wait()
	logger.LogInfo(ctx, "All goroutines completed.", nil)
}

func backgroundWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // Listen for shutdown signal
			logger.LogInfo(ctx, "Background worker shutting down...", nil)
			return
		default:
			// Simulate work
			logger.LogInfo(ctx, "Background worker is working...", nil)
			time.Sleep(2 * time.Second)
		}
	}
}
