package internal

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/controller"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/service"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router, redis redis.Redis) {
	r := repository.NewRedisRepository(redis)
	s := service.NewBasicService(r)
	basicController := controller.NewBasicController(s)

	app := router.Group("/")
	app.Get("/get", basicController.Get)
	app.Post("/set", basicController.Set)
	app.Get("/gethash/:field", basicController.GetHash)
	app.Post("/sethash", basicController.SetHash)
	app.Post("/sethashttl", basicController.SetHashTTL)
	app.Delete("/deletehash", basicController.DeleteHash)
}
