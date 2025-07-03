package user

import (
	"context"
	"github.com/ciazhar/go-start-small/internal/base_project/model/request"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
