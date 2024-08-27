package main

import (
	"context"
	"fmt"
	repository2 "github.com/ciazhar/go-zhar/examples/jwt-auth/repository"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	// Load the JWT secret key from environment variables
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	tokenRepo := repository2.NewRedisTokenRepository(rdb)

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

		getToken, err := tokenRepo.GetToken(ctx, user.ID)
		if err != nil {
			return err
		}
		if getToken != nil {
			return c.JSON(fiber.Map{"token": getToken.AuthToken})
		}

		// Generate JWT token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user.ID
		claims["exp"] = time.Now().Add(time.Hour * 5).Unix()

		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Save token in Redis
		accessToken := &repository2.AccessToken{
			UserID:      user.ID,
			AuthToken:   t,
			GeneratedAt: time.Now().Unix(),
			ExpiredAt:   time.Now().Add(time.Hour * 5).Unix(),
		}
		err = tokenRepo.SaveToken(ctx, accessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fmt.Println("token", t)

		return c.JSON(fiber.Map{"token": t})
	})

	// JWT Middleware (Custom)
	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		tokenStr := authHeader[7:]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Save token in context
		c.Locals("user", token)

		// Continue to the next handler
		return c.Next()
	})

	app.Get("/protected", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		token, err := tokenRepo.GetToken(ctx, userID)
		if err != nil || token.AuthToken != c.Get("Authorization")[7:] {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.SendString("Welcome user with ID " + string(rune(userID)))
	})

	app.Listen(":3000")
}
