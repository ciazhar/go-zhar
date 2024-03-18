package controller

import (
	"github.com/ciazhar/go-zhar/examples/message-broker/rabbitmq/basic/internal/basic/service"
	"github.com/gofiber/fiber/v2"
)

type BasicController interface {
	Publish(ctx *fiber.Ctx) error
	PublishTTL(ctx *fiber.Ctx) error
}

type basicController struct {
	aService service.BasicService
}

func (e basicController) Publish(ctx *fiber.Ctx) error {
	message := ctx.FormValue("message")
	e.aService.PublishRabbitmq(message)
	return ctx.SendStatus(fiber.StatusOK)
}

func (e basicController) PublishTTL(ctx *fiber.Ctx) error {
	message := ctx.FormValue("message")
	e.aService.PublishTTLRabbitmq(message)
	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(aService service.BasicService) BasicController {
	return basicController{
		aService: aService,
	}
}
