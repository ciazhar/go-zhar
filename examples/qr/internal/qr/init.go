package qr

import (
	"github.com/ciazhar/go-start-small/examples/qr/internal/qr/controller"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {

	c := controller.NewQrController()

	router.Get("/qr", c.GenerateQrCode)
}
