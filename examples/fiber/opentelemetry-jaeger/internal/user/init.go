package user

import (
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/controller"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/repository"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/service"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func Init(app *fiber.App, tracer trace.Tracer) {

	r := repository.NewUserRepository(tracer)
	s := service.NewUserService(r, tracer)
	c := controller.NewUserController(s, tracer)

	app.Post("/users", c.AddUser)
	app.Get("/users/:username", c.GetUserByUsername)
	app.Get("/users", c.GetAllUsers)
	app.Delete("/users/:username", c.DeleteUser)
	app.Put("/users", c.UpdateUser)
}
