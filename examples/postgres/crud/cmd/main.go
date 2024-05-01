package main

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/postgres/crud/internal"
	db "github.com/ciazhar/go-zhar/examples/postgres/crud/internal/generated/repository"
	"github.com/ciazhar/go-zhar/pkg/env"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/ciazhar/go-zhar/pkg/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {

	ctx := context.Background()

	log := logger.Init()
	env.Init("config.json", log)
	pool := postgres.Init(
		ctx,
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
	internal.Init(ctx, app, queries, pool, log)

	err := app.Listen(":" + viper.GetString("application.port"))
	if err != nil {
		return
	}
}
