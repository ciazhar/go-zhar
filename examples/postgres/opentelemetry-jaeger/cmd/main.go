package main

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/postgres/opentelemetry-jaeger/internal"
	db "github.com/ciazhar/go-zhar/examples/postgres/opentelemetry-jaeger/internal/generated/repository"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	opentelemetryjaeger "github.com/ciazhar/go-zhar/pkg/opentelemetry-jaeger"
	"github.com/ciazhar/go-zhar/pkg/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	ctx := context.Background()

	log := logger.Init(logger.Config{
		ConsoleLoggingEnabled: true,
	})
	env.Init("config.json", log)

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

	pool := postgres.Init(
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.db"),
		viper.GetString("postgres.schema"),
		viper.GetBool("application.debug"),
		log,
	)
	defer pool.Close()
	queries := db.New(pool)

	app := fiber.New()
	internal.Init(ctx, app, queries, pool, tracer, log)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		return
	}
}
