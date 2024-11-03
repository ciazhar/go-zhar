package main

import (
	"github.com/ciazhar/go-start-small/examples/api_aggregation/pkg"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/profile", func(c *fiber.Ctx) error {

		userID := c.Query("userID")
		orders := getUserProfile(userID)

		return c.JSON(orders)
	})

	app.Listen(":3003")
}

func getUserProfile(userID string) pkg.GetUserProfileResponse {
	return pkg.GetUserProfileResponse{
		UserId: userID,
		Name:   "John Doe",
		Email:  "johndoe@example.com",
	}
}
