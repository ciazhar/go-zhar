package main

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/cookie-based-auth/repository"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func main() {
	app := fiber.New()

	// Repository initialization
	authRepo := repository.NewInMemoryAuthRepository()

	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Cookie login route
	app.Post("/login", func(c *fiber.Ctx) error {
		type Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var creds Credentials
		if err := c.BodyParser(&creds); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user, err := authRepo.FindByUsername(creds.Username)
		if err != nil || user.Password != creds.Password {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Generate session ID
		sessionID := uuid.New().String()

		// Save session in Redis with expiration (e.g., 2 hours)
		err = rdb.Set(ctx, sessionID, creds.Username, time.Hour*2).Err()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Set session ID in a cookie
		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(2 * time.Hour),
			HTTPOnly: true,  // Prevents JavaScript from accessing the cookie
			Secure:   false, // Set to true in production if using HTTPS
		})

		return c.JSON(fiber.Map{"message": "Logged in successfully", "session_id": sessionID})
	})

	// Protected route that requires cookie authentication
	app.Get("/protected", cookieMiddleware, func(c *fiber.Ctx) error {
		username := c.Locals("username").(string)
		return c.SendString("Welcome " + username)
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		// Retrieve the session cookie
		sessionID := c.Cookies("session_id")

		if sessionID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Session not found")
		}

		// Delete session data from Redis
		err := rdb.Del(ctx, sessionID).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete session")
		}

		// Clear the session cookie by setting an expired cookie
		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire the cookie
			HTTPOnly: true,
		})

		return c.SendString("Logged out successfully")
	})

	app.Listen(":3000")
}

// Middleware to authenticate user based on the session ID in the cookie
func cookieMiddleware(c *fiber.Ctx) error {
	// Get the session ID from the cookie
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Fetch username from Redis using the session ID
	username, err := rdb.Get(ctx, sessionID).Result()
	if err == redis.Nil || err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Store the username in context for use in the protected route
	c.Locals("username", username)
	return c.Next()
}
