package controller

import (
	"context"
	"github.com/ciazhar/go-zhar/examples/postgres/crud/internal/model"
	"github.com/ciazhar/go-zhar/examples/postgres/crud/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProductController struct {
	ctx            context.Context
	productService service.ProductService
}

func NewProductController(ctx context.Context, productService service.ProductService) *ProductController {
	return &ProductController{ctx: ctx, productService: productService}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided name and price
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      model.CreateProductRequest  true  "Product request body"
// @Success      201
// @Failure      500      {object}  models.Error
// @Router       /products [post]
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	req := new(model.CreateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := pc.productService.CreateProduct(pc.ctx, req.Name, req.Price); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// GetProducts godoc
// @Summary      Get products
// @Description  Get a list of products with optional filtering, sorting, and pagination
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        name     query     string  false  "Product name filter"
// @Param        price    query     float64 false  "Product price filter"
// @Param        sortBy   query     string  false  "Sort by field (e.g., 'name', 'price')"
// @Param        page     query     int     false  "Page number"
// @Param        size     query     int     false  "Page size"
// @Success      200      {object}  services.GetProductsResponse
// @Failure      500      {object}  models.Error
// @Router       /products [get]
func (pc *ProductController) GetProducts(c *fiber.Ctx) error {
	name := c.Query("name")
	priceStr := c.Query("price")
	sortBy := c.Query("sortBy")
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var price float64
	if priceStr != "" {
		var err error
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid price value",
			})
		}
	}

	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid page value",
			})
		}
	}

	size := 10
	if sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid size value",
			})
		}
	}

	data, err := pc.productService.GetProducts(pc.ctx, name, price, sortBy, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

// GetProductsCursor godoc
// @Summary      Get products with cursor pagination
// @Description  Get a list of products with optional filtering and cursor-based pagination
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        name     query     string  false  "Product name filter"
// @Param        price    query     float64 false  "Product price filter"
// @Param        cursor   query     string  false  "Cursor for pagination"
// @Param        size     query     int     false  "Page size"
// @Success      200      {object}  services.GetProductsCursorResponse
// @Failure      500      {object}  models.Error
// @Router       /products/cursor [get]
func (pc *ProductController) GetProductsCursor(c *fiber.Ctx) error {
	name := c.Query("name")
	priceStr := c.Query("price")
	cursor := c.Query("cursor")
	sizeStr := c.Query("size")

	var price float64
	if priceStr != "" {
		var err error
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid price value",
			})
		}
	}

	size := 10
	if sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid size value",
			})
		}
	}

	page, err := pc.productService.GetProductsCursor(pc.ctx, name, price, cursor, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(page)
}

// UpdateProductPrice godoc
// @Summary      Update product price
// @Description  Update the price of a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "Product ID"
// @Param        product  body      services.UpdateProductPriceRequest  true  "Product request body"
// @Success      200
// @Failure      400      {object}  models.Error
// @Failure      500      {object}  models.Error
// @Router       /products/{id}/price [put]
func (pc *ProductController) UpdateProductPrice(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	req := new(model.UpdateProductPriceRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := pc.productService.UpdateProductPrice(pc.ctx, id, req.Name, req.Price); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "Product ID"
// @Success      200
// @Failure      400      {object}  models.Error
// @Failure      500      {object}  models.Error
// @Router       /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	err = pc.productService.DeleteProduct(pc.ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
