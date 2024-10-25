package main

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/qr/internal/qr"
	"github.com/ciazhar/go-start-small/pkg/config"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)



func main() {

	// Logger initialization
	logger.InitLogger(logger.LogConfig{
		LogLevel:      "debug",
		ConsoleOutput: true,
	})

	// Configuration
	config.InitConfig(
		config.Config{
			Source: "file",
			File: config.FileConfig{
				FileName: "config.json",
				FilePath: "configs",
			},
			Type: "json",
		})

	// Fiber
	app := fiber.New()

	// Module initialization
	qr.Init(app)

	// Start Fiber
	if err := app.Listen(":" + viper.GetString("application.port")); err != nil {
		logger.LogFatal(context.Background(), err, "failed to start", nil)
	}

}
