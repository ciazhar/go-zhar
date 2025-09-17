package bootstrap

import (
	ctrlUser "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/controller/rest/user"
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/ciazhar/go-zhar/pkg/middleware"
	"github.com/ciazhar/go-zhar/pkg/validator"
)

type RESTModule struct {
	v        validator.Validator
	userCtrl ctrlUser.UserController
}

func NewRESTModule(v validator.Validator, uc ctrlUser.UserController) *RESTModule {
	return &RESTModule{
		v: v, userCtrl: uc,
	}
}

// Register plugs everything into a provided *fiber.App (matches your server.NewFiberServer(func(app *fiber.App){...}))
func (m *RESTModule) Register(app *fiber.App) {
	app.Use(recover.New())
	app.Use(otelfiber.Middleware())
	app.Use(middleware.PrometheusMiddleware())

	root := app.Group("/")
	root.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
	})

	root.Get("/metrics", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Context())
		return nil
	})

	v1 := app.Group("/v1")
	v1.Post("/users", middleware.BodyParserMiddleware[request.CreateUserBodyRequest](m.v), m.userCtrl.CreateUser)
	v1.Get("/users", middleware.QueryParamParserMiddleware[request.GetUsersQueryParam](m.v), m.userCtrl.GetUsers)
	v1.Get("/users/exist", middleware.QueryParamParserMiddleware[request.UserEmailQueryParam](m.v), m.userCtrl.IsUserExistByEmail)
	v1.Get("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](m.v), m.userCtrl.GetUserByID)
	v1.Put("/users/:id",
		middleware.BodyParserMiddleware[request.UpdateUserBodyRequest](m.v),
		middleware.PathParamParserMiddleware[request.UserPathParam](m.v),
		m.userCtrl.UpdateUser,
	)
	v1.Delete("/users/:id", middleware.PathParamParserMiddleware[request.UserPathParam](m.v), m.userCtrl.DeleteUser)
	v1.Post("/users/upsert", middleware.BodyParserMiddleware[request.UpsertUserBodyRequest](m.v), m.userCtrl.UpsertUser)
}
