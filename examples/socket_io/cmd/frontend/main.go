package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html/v2" // Fiber template engine
)

func main() {
    // Set up the template engine
    engine := html.New("./web/static", ".html")

    // Initialize Fiber app with template engine
    app := fiber.New(fiber.Config{
        Views: engine,
    })

    // Route to serve the HTML template
    app.Get("/", func(c *fiber.Ctx) error {
        // Render the template and pass data to it
        return c.Render("socketio", fiber.Map{
            "Message": "Welcome to Fiber!",
        })
    })

    // Start server on port 3000
    app.Listen(":3000")
}
