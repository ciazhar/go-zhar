package internal

import (
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/controller"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/repository"
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/service"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	r := repository.NewRepository()
	s := service.NewService(r)
	c := controller.NewController(s)

	app.Get("/accident-report/:id", c.GetAccidentReportHandler)
}
