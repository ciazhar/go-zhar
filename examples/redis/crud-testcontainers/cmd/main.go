package main

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	// Logger
	log := logger.Init()

	// Environment configuration
	env.Init("config.json", log)

	// Redis configuration
	r := redis.Init(
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.password"),
		log,
	)
	defer r.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	internal.Init(app, r)

	// Start Fiber
	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
