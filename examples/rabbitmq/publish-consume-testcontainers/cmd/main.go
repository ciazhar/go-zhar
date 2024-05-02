package main

import (
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume-testcontainers/internal"
	"github.com/ciazhar/go-zhar/examples/rabbitmq/publish-consume-testcontainers/internal/model"
	"github.com/ciazhar/go-zhar/pkg/context_util"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"sync"
)

func main() {

	ctx := context_util.SetupSignalHandler()
	var wg sync.WaitGroup

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("config.json", log)

	// RabbitMQ configuration
	ra := rabbitmq.New(
		viper.GetString("application.name"),
		viper.GetString("rabbitmq.username"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
		log,
	)
	ra.CreateQueue(model.QueueBasic)
	defer ra.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	internal.Init(ctx, app, ra, &wg, log)

	// Start Fiber
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.Listen(":" + viper.GetString("application.port"))
		if err != nil {
			log.Fatalf("fiber failed to start : %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		// Once all goroutines have completed, initiate server shutdown
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}

	// Wait for all consumer goroutines to complete before exiting
	wg.Wait()

}
