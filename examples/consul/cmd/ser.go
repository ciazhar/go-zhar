package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"pkg"
	"sync"
	"syscall"
)

func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigCh:
			fmt.Printf("Received signal: %v. Initiating shutdown...\n", sig)
			cancel() // Cancel the context to stop consumer goroutines
		}
	}()

	return ctx
}

func main() {

	// Concurrent configuration
	ctx := SetupSignalHandler()
	var wg sync.WaitGroup // Use a sync.WaitGroup to manage synchronization

	// Environment configuration
	pkg.InitEnv("server.json")

	c := pkg.InitConsul(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
	)
	pkg.RetrieveConfiguration(
		c,
		viper.GetString("consul.key"),
		viper.GetString("consul.configType"),
	)
	pkg.RegisterService(
		c,
		viper.GetString("application.name"),
		viper.GetString("application.name"),
		viper.GetString("application.host"),
		viper.GetInt("application.port"))

	// Start the HTTP server
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		log.Println("Hello from server-service!")
		return c.SendStatus(fiber.StatusOK)
	})

	// Start the health check server
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.Listen(":" + viper.GetString("application.port"))
		if err != nil {
			log.Fatalf("%s: %s", "Error starting server", err)
		}

		log.Printf("Server started on port %s", viper.GetString("application.healthCheckPort"))
	}()

	select {
	case <-ctx.Done():
		// Once all goroutines have completed, initiate server shutdown
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		// Deregister the service
		pkg.DeregisterService(c, viper.GetString("application.name"))
	}
}
