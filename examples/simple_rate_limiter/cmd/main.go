package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"log"
	"time"
)

func main() {

	app := fiber.New()
	app.Use(limiter.New(
		limiter.Config{
			Max:        5,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusTooManyRequests)
			},
			LimiterMiddleware: limiter.FixedWindow{},
		},
	))

	err := app.Listen(":3001")
	if err != nil {
		log.Fatalf("fiber failed to start : %v", err)
	}
}