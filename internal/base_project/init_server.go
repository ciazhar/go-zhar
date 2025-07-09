package base_project

import (
	user3 "github.com/ciazhar/go-start-small/internal/base_project/controller/rest/user"
	"github.com/ciazhar/go-start-small/internal/base_project/repository/dummy/user"
	user2 "github.com/ciazhar/go-start-small/internal/base_project/service/user"
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitServer(fiber *fiber.App, validator validator.Validator, redisClient *redis.Client) {
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

	InitRoutes(fiber, validator, c, publicAPIRateLimiter, internalAPIRateLimiter)
}
