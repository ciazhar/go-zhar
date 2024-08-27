package main

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/paseto-auth/repository"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/o1egl/paseto"
	"github.com/spf13/viper"
)

var ctx = context.Background()

func main() {
	app := fiber.New()

	// Initialize Viper
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Load the secret key from environment variables
	secret := viper.GetString("PASETO_SECRET")
	if secret == "" {
		log.Fatal("PASETO_SECRET environment variable is not set")
	}
	fmt.Println("PASETO_SECRET:", len(secret))
	if len(secret) != 32 {
		log.Fatal("PASETO_SECRET must be exactly 32 bytes long")
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	tokenRepo := repository.NewRedisTokenRepository(rdb)

	// Initialize in-memory auth repository
	authRepo := repository.NewInMemoryAuthRepository()

	// PASETO instance
	pasetoV2 := paseto.NewV2()

	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		user, err := authRepo.FindByUsername(username)
		if err != nil || user.Password != password {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		getToken, err := tokenRepo.GetToken(ctx, user.ID)
		if err != nil {
			return err
		}
		if getToken != nil {
			return c.JSON(fiber.Map{"token": getToken.AuthToken})
		}

		// Generate PASETO token
		now := time.Now()
		exp := now.Add(5 * time.Hour)
		jsonToken := paseto.JSONToken{
			Expiration: exp,
			IssuedAt:   now,
			Subject:    fmt.Sprintf("%d", user.ID),
		}

		token, err := pasetoV2.Encrypt([]byte(secret), jsonToken, nil)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Save token in Redis
		accessToken := &repository.AccessToken{
			UserID:      user.ID,
			AuthToken:   token,
			GeneratedAt: now.Unix(),
			ExpiredAt:   exp.Unix(),
		}
		err = tokenRepo.SaveToken(ctx, accessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": token})
	})

	// PASETO Middleware (Custom)
	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokenStr := authHeader[7:]
		var jsonToken paseto.JSONToken
		err := pasetoV2.Decrypt(tokenStr, []byte(secret), &jsonToken, nil)
		if err != nil || jsonToken.Expiration.Before(time.Now()) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Save token in context
		c.Locals("user", jsonToken.Subject)

		// Continue to the next handler
		return c.Next()
	})

	// Protected route
	app.Get("/protected", func(c *fiber.Ctx) error {
		userID, _ := strconv.Atoi(c.Locals("user").(string))

		token, err := tokenRepo.GetToken(ctx, userID)
		if err != nil || token.AuthToken != c.Get("Authorization")[7:] {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.SendString("Welcome user with ID " + strconv.Itoa(userID))
	})

	app.Listen(":3000")
}
