package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal"
	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/model"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	// Create a new context for handling shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup

	// Configuration using flags for source, type, and other details
	var logLevel string
	var consoleOutput bool
	var source, configType, fileName, filePath, consulEndpoint, consulPath string

	// Parse command-line flags
	flag.StringVar(&logLevel, "log-level", "debug", "Log level (default: debug)")
	flag.BoolVar(&consoleOutput, "console-output", true, "Console output (default: true)")
	flag.StringVar(&source, "source", "file", "Configuration source (file or consul)")
	flag.StringVar(&fileName, "file-name", "config.json", "Name of the configuration file")
	flag.StringVar(&filePath, "file-path", "configs", "Path to the configuration file")
	flag.StringVar(&configType, "config-type", "json", "Configuration file type")
	flag.StringVar(&consulEndpoint, "consul-endpoint", "localhost:8500", "Consul endpoint")
	flag.StringVar(&consulPath, "consul-path", "path/to/config", "Path to the configuration in Consul")
	flag.Parse()

	// Initialize logger with parsed configuration
	logger.InitLogger(logger.LogConfig{
		LogLevel:      logLevel,
		ConsoleOutput: consoleOutput,
	})

	// Configuration using flags for source, type, and other details
	fileConfig := config.Config{
		Source: source,
		Type:   configType,
		File: config.FileConfig{
			FileName: fileName,
			FilePath: filePath,
		},
		Consul: config.ConsulConfig{
			Endpoint: consulEndpoint,
			Path:     consulPath,
		},
	}

	config.InitConfig(fileConfig)

	// RabbitMQ configuration
	ra := rabbitmq.New(
		viper.GetString("application.name"),
		viper.GetString("rabbitmq.username"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	ra.CreateQueue(model.QueueBasic)
	defer ra.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	internal.Init(ctx, app, ra, &wg)

	// Start Fiber
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.Listen(":" + viper.GetString("application.port"))
		if err != nil {
			logger.LogError(ctx, err, "Failed to start server", nil)
		}
	}()

	<-ctx.Done()
	// Once all goroutines have completed, initiate server shutdown
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.LogError(ctx, err, "Failed to shutdown server", nil)
	}

	// Wait for all consumer goroutines to complete before exiting
	wg.Wait()

}
