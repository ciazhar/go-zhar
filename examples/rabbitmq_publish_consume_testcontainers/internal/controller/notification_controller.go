package controller

import (
	"encoding/json"
	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/model"
	"github.com/ciazhar/go-start-small/examples/rabbitmq_publish_consume_testcontainers/internal/service"
	"github.com/gofiber/fiber/v2"
)

type NotificationController struct {
	notificationService *service.NotificationService
}

func (e *NotificationController) UpdateOrderStatus(ctx *fiber.Ctx) error {
	var request model.OrderStatusRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	marshal, err := json.Marshal(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	e.notificationService.PublishRabbitmq(string(marshal))
	return ctx.SendStatus(fiber.StatusOK)
}

func (e *NotificationController) SendPaymentReminder(ctx *fiber.Ctx) error {
	var request model.PaymentReminderRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	marshal, err := json.Marshal(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	e.notificationService.PublishTTLRabbitmq(string(marshal))
	return ctx.SendStatus(fiber.StatusOK)
}

func NewBasicController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}
