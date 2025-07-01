// main.go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"math/rand"
)

//func main() {
//	app := fiber.New()
//
//	// Initialize default config
//	app.Use(pprof.New())
//
//	// Or extend your config for customization
//
//	// For example, in systems where you have multiple ingress endpoints, it is common to add a URL prefix, like so:
//	app.Use(pprof.New(pprof.Config{Prefix: "/endpoint-prefix"}))
//
//	// This prefix will be added to the default path of "/debug/pprof/", for a resulting URL of: "/endpoint-prefix/debug/pprof/".
//
//	app.Get("/api/heavy", func(c *fiber.Ctx) error {
//		// Simulate CPU + memory load
//		total := 0
//		for i := 0; i < 1_000_000; i++ {
//			total += rand.Intn(100)
//		}
//		return c.SendString("Done")
//	})
//
//	log.Fatal(app.Listen(":3000"))
//}
//

func main() {
	app := fiber.New()

	// Tambahkan middleware pprof
	app.Use(pprof.New())

	// Route test
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber + pprof!")
	})

	app.Get("/api/heavy", func(c *fiber.Ctx) error {

		a()
		return c.SendString("Done")
	})

	app.Listen(":3000")
}

func a() {
	data := make([]int, 1_000_000)

	// Simulate CPU + memory load
	total := 0
	for i := 0; i < 1_000_000; i++ {
		total += rand.Intn(100)
		data[i] = total
	}
}
