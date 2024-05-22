package controller

import (
	"github.com/ciazhar/go-zhar/examples/testify-mockery/internal/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Controller struct {
	service service.ServiceInterface
}

func NewController(service service.ServiceInterface) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetAccidentReportHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID parameter is required",
		})
	}

	report, err := c.service.GetAccidentReport(ctx.UserContext(), id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get accident report",
		})
	}

	return ctx.JSON(report)
}
