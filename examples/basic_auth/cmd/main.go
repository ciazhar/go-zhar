package main

import (
	"github.com/ciazhar/go-start-small/examples/basic_auth/internal/repository"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)



func main() {
	app := fiber.New()

	// Initialize the in-memory repository
	authRepo := repository.NewInMemoryAuthRepository()

	// Apply the Basic Auth middleware
	app.Use(middleware.BasicAuthMiddleware(authRepo.FindPasswordByUsername))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the protected route!")
	})

	// Start the Fiber application on port 3000
	app.Listen(":3000")
}
