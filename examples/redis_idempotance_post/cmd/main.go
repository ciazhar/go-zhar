package main

import (
	"github.com/ciazhar/go-start-small/examples/redis_idempotance_post/internal/controller"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/products", controller.CreateProductHandler)

	app.Listen(":3000")
}
