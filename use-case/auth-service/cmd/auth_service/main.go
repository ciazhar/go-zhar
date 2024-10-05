package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/use-case/auth-service/internal/auth_service"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/env"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/postgres"
	"github.com/ciazhar/go-zhar/use-case/auth-service/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {

	env.InitEnv()
	pg := postgres.InitPostgres(context.Background())
	postgres.InitPostgresMigration()
	r := redis.InitRedis()

	app := fiber.New()
	v1 := app.Group("/api/v1")
	auth_service.Init(v1, pg, r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
