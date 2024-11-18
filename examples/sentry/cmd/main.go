package main

import (
	"github.com/ciazhar/go-start-small/examples/sentry/internal/controller"
	"github.com/ciazhar/go-start-small/examples/sentry/internal/repository"
	utils "github.com/ciazhar/go-start-small/examples/sentry/pkg"
	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func main() {

	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Environment:      "production",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Initialize Database
	db := utils.InitDB()
	defer db.Close()
	repository.SetDB(db)

	// Initialize Fiber
	app := fiber.New()

	// Middleware for Sentry
	app.Use(sentryfiber.New(sentryfiber.Options{}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Task API is running"})
	})
	controller.RegisterTaskRoutes(app)

	// Start the Fiber app
	log.Fatal(app.Listen(":3000"))
}
