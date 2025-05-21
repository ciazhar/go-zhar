package controller

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/controller/grpc/product"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	productClient product.ProductServiceClient
}

func NewOrderController(productClient product.ProductServiceClient) *OrderController {
	return &OrderController{
		productClient: productClient,
	}
}

// GetOrderDetail contoh handler REST yang ambil data produk via gRPC
func (c *OrderController) GetOrderDetail(ctx *fiber.Ctx) error {
	// misalnya ambil product ID dari query param
	productID := int32(1) // default
	if id := ctx.Query("product_id"); id != "" {
		var parsed int
		_, err := fmt.Sscanf(id, "%d", &parsed)
		if err == nil {
			productID = int32(parsed)
		}
	}

	res, err := c.productClient.GetByID(context.Background(), &product.GetByIDRequest{
		Id: productID,
	})
	if err != nil {
		log.Printf("could not get product: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get product data",
		})
	}

	// mock order response
	return ctx.JSON(fiber.Map{
		"order_id":    123,
		"product_id":  productID,
		"product":     res.Name,
		"price":       res.Price,
		"product_img": res.Image,
	})
}
