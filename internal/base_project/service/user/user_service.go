package user

import (
	"context"
	"github.com/ciazhar/go-start-small/internal/base_project/repository/postgres/user"
)

type UserService interface {
	CreateUser(ctx context.Context, requestID string) error
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}
