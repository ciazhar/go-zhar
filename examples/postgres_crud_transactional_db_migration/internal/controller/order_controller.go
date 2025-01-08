package controller

import (
	"net/http"
	"strconv"

	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	orderService service.OrderServiceInterface
}

func NewOrderController(orderService service.OrderServiceInterface) *OrderController {
	return &OrderController{orderService: orderService}
}

// PlaceOrder handles the creation of an order
func (c *OrderController) PlaceOrder(ctx *fiber.Ctx) error {
	var request struct {
		CustomerID int               `json:"customer_id"`
		Items      []model.OrderItem `json:"items"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	orderID, err := c.orderService.PlaceOrder(ctx.Context(), request.CustomerID, request.Items)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"order_id": orderID})
}

// ProcessPayment handles payment processing for an order
func (c *OrderController) ProcessPayment(ctx *fiber.Ctx) error {
	var request struct {
		OrderID int     `json:"order_id"`
		Method  string  `json:"method"`
		Amount  float64 `json:"amount"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := c.orderService.ProcessPayment(ctx.Context(), request.OrderID, request.Method, request.Amount)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "payment processed successfully"})
}

// ShipOrder handles the shipment of an order
func (c *OrderController) ShipOrder(ctx *fiber.Ctx) error {
	var request struct {
		OrderID        int    `json:"order_id"`
		TrackingNumber string `json:"tracking_number"`
		Carrier        string `json:"carrier"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	err := c.orderService.ShipOrder(ctx.Context(), request.OrderID, request.TrackingNumber, request.Carrier)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "order shipped successfully"})
}

// MarkOrderDelivered handles marking an order as delivered
func (c *OrderController) MarkOrderDelivered(ctx *fiber.Ctx) error {
	orderIDStr := ctx.Params("orderID")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid order ID"})
	}

	err = c.orderService.MarkOrderDelivered(ctx.Context(), orderID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "order marked as delivered"})
}
