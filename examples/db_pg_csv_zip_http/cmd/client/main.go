package main

import (
	"github.com/ciazhar/go-start-small/examples/db_pg_csv_zip_http/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	internal.InitClient(app)
	app.Listen(":3000")
}
