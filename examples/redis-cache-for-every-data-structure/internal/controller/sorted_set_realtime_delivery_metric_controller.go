package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/redis-cache-for-every-data-structure/internal/repository"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type DeliveryMetricController struct {
	repo *repository.SortedSetRealtimeDeliveryMetricRepository
}

func NewDeliveryMetricController(repo *repository.SortedSetRealtimeDeliveryMetricRepository) *DeliveryMetricController {
	return &DeliveryMetricController{
		repo: repo,
	}
}

func (c *DeliveryMetricController) RegisterRoutes(app *fiber.App) {
	app.Post("/add-delivery-metric", c.AddDeliveryMetricHandler)
	app.Get("/top-deliveries", c.GetTopDeliveriesHandler)
}

// AddDeliveryMetricHandler handles adding a delivery metric
func (c *DeliveryMetricController) AddDeliveryMetricHandler(ctx *fiber.Ctx) error {
	type request struct {
		OrderID   string  `json:"order_id"`
		TimeTaken float64 `json:"time_taken"`
	}

	var req request
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := c.repo.AddDeliveryMetric(context.Background(), req.OrderID, req.TimeTaken)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add delivery metric"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Delivery metric added successfully"})
}

// GetTopDeliveriesHandler handles fetching the top 5 fastest deliveries
func (c *DeliveryMetricController) GetTopDeliveriesHandler(ctx *fiber.Ctx) error {
	topDeliveries, err := c.repo.GetTopDeliveries(context.Background())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve top deliveries"})
	}

	return ctx.Status(http.StatusOK).JSON(topDeliveries)
}
