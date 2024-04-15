package qr

import (
	"github.com/ciazhar/go-zhar/examples/qr/internal/qr/controller"
	"github.com/ciazhar/go-zhar/examples/qr/internal/qr/service"
	"github.com/gofiber/fiber/v2"
)

func Init(router fiber.Router) {

	s := service.NewQrService()
	c := controller.NewQrController(s)

	router.Get("/qr", c.GenerateQrCode)
}
