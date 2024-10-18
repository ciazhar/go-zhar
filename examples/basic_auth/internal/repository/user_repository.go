package repository

import (
	"errors"
	"strings"

	"github.com/ciazhar/go-start-small/examples/basic_auth/internal/model"
	"github.com/ciazhar/go-start-small/pkg/hashing/bcrypt"
)

// AuthRepository is an interface for authentication-related operations.
type AuthRepository interface {
	FindByUsername(username string) (*model.User, error)
}

// InMemoryAuthRepository is an in-memory implementation of AuthRepository.
type InMemoryAuthRepository struct {
	users []model.User
}

// NewInMemoryAuthRepository creates a new instance of InMemoryAuthRepository with predefined users.
func NewInMemoryAuthRepository() *InMemoryAuthRepository {

	users := []model.User{
		{Username: "admin", Password: "password"},
		{Username: "user", Password: "1234"},
	}
	for i := range users {
		password , err := bcrypt.HashPassword(users[i].Password)
		if err != nil {
			panic(err)
		}
		users[i].Password = password
	}

	return &InMemoryAuthRepository{
		users: users,
	}
}

// FindPasswordByUsername searches for a user's password by their username.
func (repo *InMemoryAuthRepository) FindPasswordByUsername(username string) (string, error) {
	for _, user := range repo.users {
		if strings.EqualFold(user.Username, username) {
			return user.Password, nil
		}
	}
	return "", errors.New("user not found")
}