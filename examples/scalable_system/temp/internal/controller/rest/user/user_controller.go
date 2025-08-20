package user

import (
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/service/user"
)

type UserController struct {
	service user.UserService
}

func NewUserController(service user.UserService) UserController {
	return UserController{
		service: service,
	}
}
