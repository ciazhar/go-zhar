package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// This here will appear as a label, one can also use
	// fiberprometheus.NewWith(servicename, namespace, subsystem )
	// or
	// labels := map[string]string{"custom_label1":"custom_value1", "custom_label2":"custom_value2"}
	// fiberprometheus.NewWithLabels(labels, namespace, subsystem )
	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	//prometheus.SetSkipPaths([]string{"/ping"}) // Optional: Remove some paths from metrics
	app.Use(prometheus.Middleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Post("/some", func(c *fiber.Ctx) error {
		return c.SendString("Welcome!")
	})

	app.Listen(":4000")
}
