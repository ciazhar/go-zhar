package basic

import (
	"github.com/ciazhar/go-zhar/examples/redis/clean-architecture/internal/basic/controller"
	"github.com/ciazhar/go-zhar/examples/redis/clean-architecture/internal/basic/repository"
	"github.com/ciazhar/go-zhar/examples/redis/clean-architecture/internal/basic/service"
	"github.com/ciazhar/go-zhar/pkg/redis"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router, redis redis.Redis) {
	r := repository.NewRedisRepositoryParams(redis)
	s := service.NewBasicService(r)
	c := controller.NewBasicController(s)

	ro := router.Group("/basic")
	ro.Get("/hash/:key", c.GetBasicHash)
	ro.Post("/hash", c.SetBasicHash)
}
