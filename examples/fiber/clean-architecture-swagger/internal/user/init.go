package user

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/controller"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/repository"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {

	r := repository.NewUserRepository()
	s := service.NewUserService(r)
	c := controller.NewUserController(s)

	app.Post("/users", c.AddUser)
	app.Get("/users/:username", c.GetUserByUsername)
	app.Get("/users", c.GetAllUsers)
	app.Delete("/users/:username", c.DeleteUser)
	app.Put("/users", c.UpdateUser)
}
