package user

import (
	"github.com/ciazhar/go-start-small/examples/rest_api_service/internal/service/user"
)

type UserController struct {
	service user.UserService
}

func NewUserController(service user.UserService) UserController {
	return UserController{
		service: service,
	}
}
