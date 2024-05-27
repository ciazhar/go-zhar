package main

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	opentelemetryjaeger "github.com/ciazhar/go-zhar/pkg/opentelemetry-jaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})

	env.Init("config-user.json", log)

	tp, tracer := opentelemetryjaeger.InitTracer(
		viper.GetString("application.name"),
		viper.GetString("opentelemetry.url"),
		log,
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Infof("Error shutting down tracer provider: %v", err)
		}
	}()

	app := fiber.New()

	user.Init(app, tracer)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}
