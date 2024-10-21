package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/ciazhar/go-start-small/pkg/hashing/bcrypt"
	"github.com/gofiber/fiber/v2"
)

// BasicAuthMiddleware Basic Auth middleware for password hashed with bcrypt
func BasicAuthMiddleware(isUserExist func(username string) (string, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const basicPrefix = "Basic "
		authHeader := c.Get("Authorization")

		// Return 401 if no Authorization header is provided
		if authHeader == "" {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Missing Authorization header")
		}

		// Check if the Authorization header starts with "Basic "
		if len(authHeader) <= len(basicPrefix) || !strings.HasPrefix(strings.ToUpper(authHeader[:len(basicPrefix)]), basicPrefix) {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid Authorization header")
		}

		// Decode the Base64-encoded username:password
		encodedCredentials := authHeader[len(basicPrefix):]
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Malformed base64 credentials")
		}

		// Split the decoded credentials into username and password
		credentials := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(credentials) != 2 {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Malformed credentials")
		}
		username, providedPassword := credentials[0], credentials[1]

		// Check if the user exists in the database
		storedPasswordHash, err := isUserExist(username)
		if err != nil {
			// Log the error for debugging, but don't leak details to the client
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid username or password")
		}

		// Compare the provided password with the stored bcrypt-hashed password
		if bcrypt.ValidatePassword(providedPassword, storedPasswordHash) {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid username or password")
		}

		// Authentication successful, proceed to the next handler
		return c.Next()
	}
}
