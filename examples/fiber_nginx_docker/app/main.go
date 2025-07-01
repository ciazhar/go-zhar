package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func main() {
	app := fiber.New()

	port := os.Getenv("PORT")
	instance := os.Getenv("INSTANCE")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Halo dari %s di port %s 🚀", instance, port))
	})

	app.Listen(":" + port)
}
