package main

import (
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	internal.Init(app)
	app.Listen(":3000")
}
