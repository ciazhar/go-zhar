package user

import (
	"github.com/ciazhar/go-zhar/examples/rest_api_service/internal/service/user"
)

type UserController struct {
	service user.UserService
}

func NewUserController(service user.UserService) UserController {
	return UserController{
		service: service,
	}
}
