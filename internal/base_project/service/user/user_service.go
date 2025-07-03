package user

import (
	"context"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
	"github.com/ciazhar/go-start-small/internal/base_project/repository/postgres/user"
)

type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}
