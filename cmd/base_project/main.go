package main

import (
	"github.com/ciazhar/go-start-small/internal/base_project"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	logger.InitLogger(logger.LogConfig{
		LogLevel:      "debug",
		LogFile:       "./logfile.log",
		MaxSize:       10,
		MaxBackups:    5,
		MaxAge:        30,
		Compress:      true,
		ConsoleOutput: true,
	})

	v := validator.New("id")

	f := fiber.New()
	f.Use(recover.New())
	f.Use(middleware.RequestID())
	f.Use(middleware.Logger())

	base_project.InitServer(f, v)

	f.Listen(":3000")
}
