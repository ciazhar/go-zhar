package controller

import (
	"github.com/ciazhar/go-zhar/examples/qr/internal/qr/service"
	"github.com/gofiber/fiber/v2"
	"os"
)

type QrController struct {
	qrService *service.QrService
}

func (q *QrController) GenerateQrCode(ctx *fiber.Ctx) error {

	url := ctx.Query("url")
	dimension := ctx.QueryInt("dimension", 255)
	base64Image := ctx.Query("base64Image", "")

	fileName, err := q.qrService.GenerateQrCode(url, dimension, base64Image)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = ctx.Status(fiber.StatusOK).SendFile(fileName)
	if err != nil {
		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	return nil
}

func NewQrController(qrService *service.QrService) *QrController {
	return &QrController{
		qrService: qrService,
	}
}
