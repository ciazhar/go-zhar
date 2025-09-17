package user

import "github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/repository/postgres/user"

type userService struct {
	repo user.UserRepositoryContract
}

func NewUserService(repo user.UserRepositoryContract) UserService {
	return &userService{repo: repo}
}
