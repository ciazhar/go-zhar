package main

import (
	"github.com/ciazhar/go-zhar/examples/cache/redis/basic/internal/basic"
	"github.com/ciazhar/go-zhar/pkg/cache/redis"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func main() {

	// Environment configuration
	env.Init("config.json")

	// Redis configuration
	r := redis.Init(viper.GetString("redis.host"), viper.GetInt("redis.port"), viper.GetString("redis.password"))
	defer r.Close()

	// Fiber configuration
	app := fiber.New()

	// Module initialization
	basic.Init(app, r)

	// Start Fiber
	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatal(err)
	}
}
