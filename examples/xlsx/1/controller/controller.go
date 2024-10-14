package controller

import (
	"github.com/ciazhar/go-zhar/examples/xlsx/1/service"
	"github.com/gofiber/fiber/v2"
)

// DownloadXLSX handles the HTTP request and sends the generated XLSX file as a response.
func DownloadXLSX(c *fiber.Ctx) error {
	xlsxData, err := service.GenerateXLSX()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate XLSX file")
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\"report.xlsx\"")
	return c.Send(xlsxData)
}
