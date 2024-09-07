package main

import (
	repository2 "github.com/ciazhar/go-zhar/examples/api-key-auth/repository"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func main() {
	app := fiber.New()

	// Initialize in-memory auth repository
	authRepo := repository2.NewInMemoryAuthRepository()

	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		user, err := authRepo.FindByUsername(username)
		if err != nil || user.Password != password {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.JSON(fiber.Map{"token": user.ApiKey})
	})

	// API Key Middleware (Custom)
	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		apiKey := authHeader[7:]

		// Check if the API key exists in Redis
		token, err := authRepo.FindByApiKey(apiKey)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Save user ID in context
		c.Locals("user_id", token.ID)

		// Continue to the next handler
		return c.Next()
	})

	app.Get("/protected", func(c *fiber.Ctx) error {
		// Get user ID from context
		userID := c.Locals("user_id")

		return c.SendString("Welcome user with ID " + strconv.Itoa(userID.(int)))
	})

	app.Listen(":3000")
}
