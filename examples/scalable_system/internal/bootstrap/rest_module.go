package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/compress"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	redisv9 "github.com/redis/go-redis/v9"

	ctrlUser "github.com/ciazhar/go-zhar/examples/scalable_system/internal/controller/rest/user"
	"github.com/ciazhar/go-zhar/examples/scalable_system/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/middleware"
	"github.com/ciazhar/go-zhar/pkg/rate_limiter"
	"github.com/ciazhar/go-zhar/pkg/validator"
)

type RESTModule struct {
	v          validator.Validator
	rdb        *redisv9.Client
	userCtrl   ctrlUser.UserController
	publicRL   rate_limiter.RateLimiter
	internalRL rate_limiter.RateLimiter
}

func NewRESTModule(v validator.Validator, rdb *redisv9.Client, uc ctrlUser.UserController) *RESTModule {
	//store := rate_limiter.NewInMemoryStore(1*time.Minute, 1*time.Minute)
	store := rate_limiter.NewRedisStore(rdb)

	public := rate_limiter.NewRateLimiter(rate_limiter.RateLimitConfig{
		Type:   rate_limiter.FixedWindowType,
		Key:    rate_limiter.ApiKey,
		Store:  store,
		Limit:  10,
		Window: 1 * time.Minute,
	})
	internal := rate_limiter.NewRateLimiter(rate_limiter.RateLimitConfig{
		Type:   rate_limiter.LeakyBucketType,
		Key:    rate_limiter.IpAddress,
		Store:  store,
		Limit:  1000,
		Window: 1 * time.Minute,
	})

	return &RESTModule{
		v: v, rdb: rdb, userCtrl: uc,
		publicRL: public, internalRL: internal,
	}
}

// Register plugs everything into a provided *fiber.App (matches your server.NewFiberServer(func(app *fiber.App){...}))
func (m *RESTModule) Register(app *fiber.App) {
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger())

	compressionMiddleware := compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	})

	root := app.Group("/", middleware.RateLimitMiddleware(m.internalRL))
	root.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
	})

	v1 := app.Group("/v1", middleware.RateLimitMiddleware(m.publicRL))
	v1.Post("/users", middleware.BodyParserMiddleware[request.CreateUserBodyRequest](m.v), m.userCtrl.CreateUser)
	v1.Get("/users", middleware.QueryParamParserMiddleware[request.GetUsersQueryParam](m.v), m.userCtrl.GetUsers)
	v1.Get("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](m.v), m.userCtrl.GetUserByID)
	v1.Put("/users/:id",
		middleware.BodyParserMiddleware[request.UpdateUserBodyRequest](m.v),
		middleware.PathParamParserMiddleware[request.UserPathParam](m.v),
		m.userCtrl.UpdateUser,
	)
	v1.Delete("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](m.v), m.userCtrl.DeleteUser)

	v1.Get("/big-json-compressed", compressionMiddleware, func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"users": []fiber.Map{
				{
					"id":    1,
					"name":  "Alice Johnson",
					"email": "alice@example.com",
					"roles": []string{"admin", "editor"},
					"profile": fiber.Map{
						"age":     29,
						"country": "USA",
						"preferences": fiber.Map{
							"language": "en",
							"timezone": "UTC-5",
						},
					},
				},
				{
					"id":    2,
					"name":  "Bob Smith",
					"email": "bob@example.com",
					"roles": []string{"viewer"},
					"profile": fiber.Map{
						"age":     34,
						"country": "UK",
						"preferences": fiber.Map{
							"language": "en",
							"timezone": "UTC+0",
						},
					},
				},
			},
			"products": []fiber.Map{
				{
					"id":    101,
					"name":  "Laptop Pro 15",
					"price": 1599.99,
					"stock": 42,
					"categories": []string{
						"electronics",
						"computers",
					},
					"attributes": fiber.Map{
						"cpu":     "Intel i7",
						"ram":     "16GB",
						"storage": "512GB SSD",
					},
				},
				{
					"id":    102,
					"name":  "Wireless Headphones",
					"price": 199.99,
					"stock": 350,
					"categories": []string{
						"electronics",
						"audio",
					},
					"attributes": fiber.Map{
						"battery_life":     "30h",
						"noise_cancelling": true,
					},
				},
			},
			"transactions": []fiber.Map{
				{
					"id":         "T0001",
					"user_id":    1,
					"product_id": 101,
					"quantity":   1,
					"total":      1599.99,
					"status":     "completed",
					"timestamp":  "2025-08-19T10:15:00Z",
				},
				{
					"id":         "T0002",
					"user_id":    2,
					"product_id": 102,
					"quantity":   2,
					"total":      399.98,
					"status":     "pending",
					"timestamp":  "2025-08-19T11:20:00Z",
				},
			},
		})
	})
}
