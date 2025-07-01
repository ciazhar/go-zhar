package main

import (
	"github.com/ciazhar/go-start-small/examples/db_ch_csv_zip_http/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100 MB
	})
	internal.InitClient(app)
	app.Listen(":3000")
}
