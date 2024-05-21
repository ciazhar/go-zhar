package main

import (
	"github.com/ciazhar/go-zhar/pkg/consul"
	"github.com/ciazhar/go-zhar/pkg/context_util"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"sync"
)

func main() {

	// Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	// Concurrent configuration
	ctx := context_util.SetupSignalHandler()
	var wg sync.WaitGroup // Use a sync.WaitGroup to manage synchronization

	// Environment configuration
	env.Init("server.json", log)
	c := consul.Init(
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
		viper.GetString("consul.scheme"),
		log,
	)
	c.RetrieveConfiguration(viper.GetString("consul.key"), viper.GetString("consul.configType"))
	c.RegisterService(
		viper.GetString("application.name"),
		viper.GetString("application.name"),
		viper.GetString("application.host"),
		viper.GetInt("application.port"),
	)

	// Start the HTTP server
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		log.Infof("Hello from my-service!")
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

		log.Infof("Server started on port %s", viper.GetString("application.healthCheckPort"))
	}()

	select {
	case <-ctx.Done():
		// Once all goroutines have completed, initiate server shutdown
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}

		// Deregister the service
		c.DeregisterService(viper.GetString("application.name"))
	}

	// Wait for all consumer goroutines to complete before exiting
	wg.Wait()
}
