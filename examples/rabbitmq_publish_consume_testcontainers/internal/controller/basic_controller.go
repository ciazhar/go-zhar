package controller

import (
	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/service"
	"github.com/gofiber/fiber/v2"
)

type BasicController struct {
	aService *service.BasicService
}

func (e *BasicController) Publish(ctx *fiber.Ctx) error {
	message := ctx.FormValue("message")
	e.aService.PublishRabbitmq(message)
	return ctx.SendStatus(fiber.StatusOK)
}

func (e *BasicController) PublishTTL(ctx *fiber.Ctx) error {
	message := ctx.FormValue("message")
	e.aService.PublishTTLRabbitmq(message)
	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(aService *service.BasicService) *BasicController {
	return &BasicController{
		aService: aService,
	}
}
