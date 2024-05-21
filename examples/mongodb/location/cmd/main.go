package main

import (
	"github.com/ciazhar/go-zhar/examples/mongodb/location/internal/location"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/mongo"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	// Environment configuration
	env.Init("config.json", log)

	// Mongo configuration
	mongoDatabase := mongo.Init(
		viper.GetString("mongo.hosts"),
		viper.GetString("mongo.username"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.database"),
		log,
	)

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	location.Init(app, mongoDatabase)

	// Start Fiber
	if err := app.Listen(":" + viper.GetString("application.port")); err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
