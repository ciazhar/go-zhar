package bootstrap

import (
	user3 "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/controller/rest/user"
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/repository/dummy/user"
	user2 "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/service/user"
	"github.com/ciazhar/go-zhar/pkg/middleware"
	"github.com/ciazhar/go-zhar/pkg/rate_limiter"
	"github.com/ciazhar/go-zhar/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitServer(fiberApp *fiber.App, validator validator.Validator, redisClient *redis.Client) {

	fiberApp.Use(recover.New())
	fiberApp.Use(middleware.RequestID())
	fiberApp.Use(middleware.Logger())

	r := user.NewUserRepository()
	s := user2.NewUserService(r)
	c := user3.NewUserController(s)

	store := rate_limiter.NewRedisStore(redisClient)
	//store := rate_limiter.NewInMemoryStore(5*time.Minute, 10*time.Minute)
	publicAPIRateLimiter := rate_limiter.NewRateLimiter(rate_limiter.RateLimitConfig{
		Type:   rate_limiter.FixedWindowType,
		Key:    rate_limiter.ApiKey,
		Store:  store,
		Limit:  10,
		Window: 1 * time.Minute,
	})

	internalAPIRateLimiter := rate_limiter.NewRateLimiter(rate_limiter.RateLimitConfig{
		Type:   rate_limiter.LeakyBucketType,
		Key:    rate_limiter.IpAddress,
		Store:  store,
		Limit:  1000,
		Window: 1 * time.Minute,
	})

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
