package main

import (
	"github.com/ciazhar/go-zhar/examples/qr/internal/qr"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("config.json", log)

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	qr.Init(app)

	// Start Fiber
	if err := app.Listen(":" + viper.GetString("application.port")); err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
