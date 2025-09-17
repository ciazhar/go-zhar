package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/temp/internal/model/response"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
	GetUserByID(ctx context.Context, id string) (*response.User, error)
	UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, page, limit int) ([]response.User, int64, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
