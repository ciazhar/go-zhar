package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/internal/model/response"
	"github.com/ciazhar/go-zhar/examples/scalable_system/internal/repository/dummy/user"
)

type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
	GetUserByID(ctx context.Context, id string) (*response.User, error)
	UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context, req request.GetUsersQueryParam) ([]response.User, int64, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}
