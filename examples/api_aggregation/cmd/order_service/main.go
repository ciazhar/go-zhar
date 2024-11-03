package main

import (
	"github.com/ciazhar/go-start-small/examples/api_aggregation/pkg"
	"github.com/gofiber/fiber/v2"
)

func main(){
	app := fiber.New()

	app.Get("/orders", func(c *fiber.Ctx) error {

		userID := c.Query("userID")
		orders := getUserOrders(userID)

		return c.JSON(orders)
	})

	app.Listen(":3001")
}

func getUserOrders(userID string) []pkg.GetUserOrdersResponse {
	return []pkg.GetUserOrdersResponse{
		{
			OrderID: "123",
			Items:   []string{"Laptop", "Mouse"},
			Status:  "Delivered",	
		},
		{
			OrderID: "124",
			Items:   []string{"Keyboard", "Monitor"},
			Status:  "Processing",
		},
	}
}