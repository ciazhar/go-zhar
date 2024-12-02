package controller

import (
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	orderService service.OrderServiceInterface
}

func NewOrderController(orderService service.OrderServiceInterface) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (uc *OrderController) CreateOrder(c *fiber.Ctx) error {

	var request model.OrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := uc.orderService.CreateOrder(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (uc *OrderController) GetAllOrders(c *fiber.Ctx) error {
	orders, err := uc.orderService.GetAllOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (uc *OrderController) GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	order, err := uc.orderService.GetOrder(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}

func (uc *OrderController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	err = uc.orderService.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
