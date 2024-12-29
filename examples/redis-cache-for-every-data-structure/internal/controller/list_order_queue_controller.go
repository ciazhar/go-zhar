package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis-cache-for-every-data-structure/internal/repository"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type OrderQueueController struct {
	repo *repository.ListOrderQueueRepository
}

// NewOrderQueueController initializes the controller with the repository
func NewOrderQueueController(repo *repository.ListOrderQueueRepository) *OrderQueueController {
	return &OrderQueueController{
		repo: repo,
	}
}

// RegisterRoutes registers the routes for the OrderQueue controller
func (c *OrderQueueController) RegisterRoutes(app *fiber.App) {
	app.Post("/orders/queue", c.AddOrderToQueueHandler)
	app.Post("/orders/next", c.ProcessNextOrderHandler)
}

// AddOrderToQueueHandler handles adding an order to the queue
func (c *OrderQueueController) AddOrderToQueueHandler(ctx *fiber.Ctx) error {
	var request struct {
		OrderID string `json:"order_id"`
	}

	if err := ctx.BodyParser(&request); err != nil || request.OrderID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body or missing order_id")
	}

	if err := c.repo.AddOrderToQueue(context.Background(), request.OrderID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add order to queue")
	}

	return ctx.Status(fiber.StatusOK).SendString("Order added to queue successfully")
}

// ProcessNextOrderHandler handles processing the next order in the queue
func (c *OrderQueueController) ProcessNextOrderHandler(ctx *fiber.Ctx) error {
	orderID, err := c.repo.ProcessNextOrder(context.Background())
	if err != nil {
		if err == redis.Nil {
			return fiber.NewError(fiber.StatusNotFound, "No orders in queue")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to process next order")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"processed_order_id": orderID,
	})
}
