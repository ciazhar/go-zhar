package controller

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	Service *service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{Service: service}
}

func (c *ProductController) GetProducts(ctx *fiber.Ctx) error {
	products, err := c.Service.GetProducts(context.Background())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(products)
}

func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	type CreateProductRequest struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Stock int     `json:"stock"`
	}

	var req CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	id, err := c.Service.CreateProduct(context.Background(), req.Name, req.Price, req.Stock)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
