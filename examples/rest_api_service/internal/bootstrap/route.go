package bootstrap

import (
	"github.com/ciazhar/go-start-small/examples/rest_api_service/internal/controller/rest/user"
	"github.com/ciazhar/go-start-small/examples/rest_api_service/internal/model/request"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(fiberApp *fiber.App, validator validator.Validator, c user.UserController, publicAPIRateLimiter rate_limiter.RateLimiter, internalAPIRateLimiter rate_limiter.RateLimiter) {
	root := fiberApp.Group("/", middleware.RateLimitMiddleware(internalAPIRateLimiter))
	root.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "ok",
		})
	})

	v1 := fiberApp.Group("/v1", middleware.RateLimitMiddleware(publicAPIRateLimiter))
	v1.Post("/users", middleware.BodyParserMiddleware[request.CreateUserBodyRequest](validator), c.CreateUser)
	v1.Get("/users", middleware.QueryParamParserMiddleware[request.GetUsersQueryParam](validator), c.GetUsers)
	v1.Get("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.GetUserByID)
	v1.Put("/users/:id", middleware.BodyParserMiddleware[request.UpdateUserBodyRequest](validator), middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.UpdateUser)
	v1.Delete("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](validator), c.DeleteUser)
}
