package main

import (
	"github.com/ciazhar/go-start-small/internal/base_project"
	"github.com/ciazhar/go-start-small/pkg/logger"
	"github.com/ciazhar/go-start-small/pkg/middleware"
	"github.com/ciazhar/go-start-small/pkg/rate_limiter"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"
)

func main() {

	logger.InitLogger(logger.LogConfig{
		LogLevel:      "debug",
		LogFile:       "./logfile.log",
		MaxSize:       10,
		MaxBackups:    5,
		MaxAge:        30,
		Compress:      true,
		ConsoleOutput: true,
	})

	v := validator.New("id")

	//redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	//store := rate_limiter.NewRedisStore(redisClient)
	store := rate_limiter.NewInMemoryStore(5*time.Minute, 10*time.Minute)
	limiter := rate_limiter.NewRateLimiter(rate_limiter.RateLimitConfig{
		Type:   rate_limiter.FixedWindowType,
		Key:    rate_limiter.ApiKey,
		Store:  store,
		Limit:  10,
		Window: 1 * time.Minute,
	})

	f := fiber.New()
	f.Use(recover.New())
	f.Use(middleware.RequestID())
	f.Use(middleware.Logger())
	f.Use(middleware.RateLimitMiddleware(limiter))

	base_project.InitServer(f, v)

	f.Listen(":3000")
}
