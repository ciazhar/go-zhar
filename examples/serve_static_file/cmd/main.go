package main

import (
	"log"
	"github.com/ciazhar/go-start-small/examples/serve_static_file/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Serve the index.html file for the root route
	app.Get("/", func(c *fiber.Ctx) error {

		html, err := web.EmbedFs.ReadFile("static/index.html")
		if err != nil {
			return c.SendStatus(500)
		}
		c.Set("Content-Type", "text/html")

		return c.SendString(string(html))
	})

	// Start the server
	log.Fatal(app.Listen(":3000"))
}