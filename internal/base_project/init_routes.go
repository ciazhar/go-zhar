package base_project

import (
	user3 "github.com/ciazhar/go-start-small/internal/base_project/controller/rest/user"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(fiber *fiber.App, c user3.UserController) {

	v1 := fiber.Group("/v1")
	v1.Post("/user", c.CreateUser)
}
