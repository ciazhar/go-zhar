package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type UserProfile struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

var inMemoryDB = map[string]*UserProfile{}

func GetUserProfileFromDB(userID string) (*UserProfile, error) {
	if profile, ok := inMemoryDB[userID]; ok {
		return profile, nil
	}
	// default fallback
	return &UserProfile{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
		Photo: "avatar.jpg",
	}, nil
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6377",
	})

	r := fiber.New()

	r.Get("/profile/:id", func(c *fiber.Ctx) error {
		ctx := context.Background()
		userID := c.Params("id")
		key := fmt.Sprintf("user:profile:%s", userID)

		// Try to get from Redis
		val, err := rdb.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			// Cache miss → get from DB
			profile, err := GetUserProfileFromDB(userID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user from DB")
			}

			// Cache it to Redis
			jsonVal, _ := json.Marshal(profile)
			rdb.Set(ctx, key, jsonVal, 10*time.Minute)

			return c.JSON(profile)
		} else if err != nil {
			log.Printf("Redis error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Redis error")
		}

		// Cache hit
		var profile UserProfile
		if err := json.Unmarshal([]byte(val), &profile); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse cache")
		}

		return c.JSON(profile)
	})

	r.Put("/profile/:id", func(c *fiber.Ctx) error {
		ctx := context.Background()
		userID := c.Params("id")
		key := fmt.Sprintf("user:profile:%s", userID)

		var updatedProfile UserProfile
		if err := c.BodyParser(&updatedProfile); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		// Simulate DB update
		inMemoryDB[userID] = &updatedProfile
		log.Printf("User %s updated in inMemoryDB", userID)

		// Invalidate cache
		rdb.Del(ctx, key)

		return c.JSON(fiber.Map{
			"message": "User updated and cache cleared",
		})
	})

	log.Fatal(r.Listen(":3000"))
}
