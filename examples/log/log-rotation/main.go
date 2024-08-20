package main

import (
	"flag"
	"github.com/ciazhar/go-zhar/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var logFileDir = flag.String("log", logger.LogFileDir, "log file directory")
	flag.Parse()

	log := logger.Init(logger.Config{
		FileLoggingEnabled: true,
		Directory:          *logFileDir,
		Filename:           "service-xyz.log",
		MaxSize:            100, // megabytes
		MaxBackups:         3,
		MaxAge:             28, // days
	})

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		log.Infof("Hello, World!")
		return c.SendString("Hello, World!")
	})

	app.Listen(":3002")

}
