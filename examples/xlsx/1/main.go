package main

import (
	"github.com/ciazhar/go-zhar/examples/xlsx/1/controller"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/download-xlsx", controller.DownloadXLSX)

	app.Listen(":8081")
}
