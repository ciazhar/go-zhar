package main

import (
	"encoding/base64"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
}

// AuthRepository is an interface for authentication-related operations.
type AuthRepository interface {
	FindByUsername(username string) (*User, error)
}

// InMemoryAuthRepository is an in-memory implementation of AuthRepository.
type InMemoryAuthRepository struct {
	users []User
}

// NewInMemoryAuthRepository creates a new instance of InMemoryAuthRepository with predefined users.
func NewInMemoryAuthRepository() *InMemoryAuthRepository {
	return &InMemoryAuthRepository{
		users: []User{
			{Username: "admin", Password: "password"},
			{Username: "user", Password: "1234"},
		},
	}
}

// FindByUsername searches for a user by their username.
func (repo *InMemoryAuthRepository) FindByUsername(username string) (*User, error) {
	for _, user := range repo.users {
		if strings.EqualFold(user.Username, username) {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func basicAuthMiddleware(repo AuthRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			c.Set("WWW-Authenticate", `Basic realm="restricted"`)
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		if len(authHeader) > 6 && strings.ToUpper(authHeader[:6]) == "BASIC " {
			decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) == 2 {
				user, err := repo.FindByUsername(credentials[0])
				if err == nil && user.Password == credentials[1] {
					return c.Next() // Authorized, proceed to the next handler
				}
			}
		}

		c.Set("WWW-Authenticate", `Basic realm="restricted"`)
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
}

func main() {
	app := fiber.New()

	// Initialize the in-memory repository
	authRepo := NewInMemoryAuthRepository()

	// Apply the Basic Auth middleware
	app.Use(basicAuthMiddleware(authRepo))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the protected route!")
	})

	// Start the Fiber application on port 3000
	app.Listen(":3000")
}
