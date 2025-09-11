package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
)

type UserRepositoryContract interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
	GetUserByID(ctx context.Context, id string) (*response.User, error)
	UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, page, limit int) ([]response.User, int64, error)
}
