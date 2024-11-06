package controller

import (
	"github.com/ciazhar/go-start-small/examples/mock_testing/internal/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type AccidentReportController struct {
	service service.AccidentReportServiceInterface
}

func NewAccidentReportController(service service.AccidentReportServiceInterface) *AccidentReportController {
	return &AccidentReportController{service: service}
}

func (c *AccidentReportController) GetAccidentReportHandler(ctx *fiber.Ctx) error {
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
