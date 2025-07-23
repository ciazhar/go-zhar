package bootstrap

import (
	user3 "github.com/ciazhar/go-start-small/examples/rest_api_service/internal/controller/rest/user"
	"github.com/ciazhar/go-start-small/examples/rest_api_service/internal/repository/dummy/user"
	user2 "github.com/ciazhar/go-start-small/examples/rest_api_service/internal/service/user"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/ciazhar/go-start-small/pkg/validator"
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

	InitRoutes(fiberApp, validator, c, publicAPIRateLimiter, internalAPIRateLimiter)
}
