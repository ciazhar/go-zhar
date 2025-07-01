package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, requestID string) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}
