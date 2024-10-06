package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIDMiddleware(c *fiber.Ctx) error {
	// Generate a new requestID using UUID
	requestID := uuid.New().String()

	// Set the requestID in the response headers (optional)
	c.Set("X-Request-ID", requestID)

	// Store requestID in the context
	ctx := c.Context()
	ctx.SetUserValue("requestID", requestID)

	// Call the next handler in the chain with the updated context
	return c.Next()
}
