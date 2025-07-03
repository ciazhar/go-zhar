package base_project

import (
	user3 "github.com/ciazhar/go-start-small/internal/base_project/controller/rest/user"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(fiber *fiber.App, validator validator.Validator, c user3.UserController) {

	v1 := fiber.Group("/v1")
	v1.Post("/users", middleware.BodyParserMiddleware[request.CreateUserBodyRequest](validator), c.CreateUser)
	v1.Get("/users", middleware.QueryParamParserMiddleware[request.GetUsersQueryParam](validator), c.GetUsers)
	v1.Get("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.GetUserByID)
	v1.Put("/users/:id", middleware.BodyParserMiddleware[request.UpdateUserBodyRequest](validator), middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.UpdateUser)
	v1.Delete("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.DeleteUser)
}
