package main

import (
	"github.com/ciazhar/go-start-small/examples/api_aggregation/pkg"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/recommendations", func(c *fiber.Ctx) error {

		userID := c.Query("userID")
		orders := getProductRecommendations(userID)

		return c.JSON(orders)
	})

	app.Listen(":3002")
}

func getProductRecommendations(userID string) []pkg.GetProductRecommendationsResponse {
	return []pkg.GetProductRecommendationsResponse{
		{
			ProductId: "123",
			Name:      "Laptop",
			Price:     999.99,
		},
		{
			ProductId: "321",
			Name:      "Wireless Headphones",
			Price:     79.99,
		},
		{
			ProductId: "322",
			Name:      "Smartphone Stand",
			Price:     19.99,
		},
	}
}
