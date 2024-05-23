package main

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture/internal/user"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})
	env.Init("config.json", log)

	app := fiber.New()

	user.Init(app)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
