package base_project

import (
	user3 "github.com/ciazhar/go-start-small/internal/base_project/controller/rest/user"
	"github.com/ciazhar/go-start-small/internal/base_project/repository/postgres/user"
	user2 "github.com/ciazhar/go-start-small/internal/base_project/service/user"
	"github.com/ciazhar/go-start-small/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func InitServer(fiber *fiber.App, validator validator.Validator) {
	r := user.NewUserRepository()
	s := user2.NewUserService(r)
	c := user3.NewUserController(s)

	InitRoutes(fiber, validator, c)
}
