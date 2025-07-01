package main

import (
	"github.com/ciazhar/go-start-small/internal/base_project"
	"github.com/gofiber/fiber/v2"
)

func main() {

	f := fiber.New()

	base_project.InitServer(f)

	f.Listen(":3000")
}
