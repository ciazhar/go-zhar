package bootstrap

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	redisv9 "github.com/redis/go-redis/v9"

	ctrlUser "github.com/ciazhar/go-zhar/examples/rest_api_service/internal/controller/rest/user"
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/model/request"
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
}
