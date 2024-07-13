package internal

import (
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/controller"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/repository"
	"github.com/ciazhar/go-zhar/examples/redis/crud-testcontainers/internal/service"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router, redis *redis.Redis) {
	r := repository.NewRedisRepository(redis)
	s := service.NewBasicService(r)
	basicController := controller.NewBasicController(s)

	app := router.Group("/")
	app.Get("/string", basicController.Get)
	app.Post("/string", basicController.Set)
	app.Delete("/string", basicController.Delete)

	app.Get("/hash/:field", basicController.GetHash)
	app.Post("/hash", basicController.SetHash)
	app.Post("/hash-ttl", basicController.SetHashTTL)
	app.Delete("/hash", basicController.DeleteHash)

	app.Get("/list", basicController.GetList)
	app.Post("/list", basicController.SetList)
	app.Delete("/list", basicController.DeleteList)
}
