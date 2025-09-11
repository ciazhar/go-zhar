package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
)

type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error
	GetUserByID(ctx context.Context, id string) (*response.User, error)
	GetUsersWithPagination(ctx context.Context, page, size int) ([]response.User, int64, error)
	IsUserExistsByEmail(ctx context.Context, email string) (bool, error)
	SoftDeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, req request.UpdateUserBodyRequest) error
	UpsertUserByID(ctx context.Context, req request.UpsertUserBodyRequest) error
}
