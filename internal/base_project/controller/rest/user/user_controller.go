package user

import (
	"github.com/ciazhar/go-start-small/internal/base_project/service/user"
)

type UserController struct {
	service user.UserService
}

func NewUserController(service user.UserService) UserController {
	return UserController{
		service: service,
	}
}
