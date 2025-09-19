package user

import (
	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/service/user"
)

type UserController struct {
	service user.UserService
}

func NewUserController(service user.UserService) UserController {
	return UserController{
		service: service,
	}
}
