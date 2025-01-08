package controller

import (
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	Service *service.CustomerService
}

func NewCustomerController(service *service.CustomerService) *CustomerController {
	return &CustomerController{Service: service}
}

func (c *CustomerController) CreateCustomer(ctx *fiber.Ctx) error {
	var req model.Customer

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	id, err := c.Service.CreateCustomer(ctx.Context(), req.Name, req.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
