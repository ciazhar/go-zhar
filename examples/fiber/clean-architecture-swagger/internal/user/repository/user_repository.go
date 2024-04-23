package repository

import (
	"fmt"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/model"
)

type UserRepository interface {
	AddUser(user model.User)
	GetUserByUsername(username string) (*model.User, error)
	GetAllUsers() map[string]model.User
	DeleteUser(username string) error
	UpdateUser(user model.User) error
}

// userRepository represents a repository for managing users
type userRepository struct {
	users map[string]model.User
}

// AddUser adds a new user to the repository
func (r *userRepository) AddUser(user model.User) {
	r.users[user.Username] = user
}

// GetUserByUsername retrieves a user by their username
func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetAllUsers retrieves all users from the repository
func (r *userRepository) GetAllUsers() map[string]model.User {
	return r.users
}

// DeleteUser deletes a user from the repository
func (r *userRepository) DeleteUser(username string) error {
	_, ok := r.users[username]
	if !ok {
		return fmt.Errorf("user not found for username: %s", username)
	}

	delete(r.users, username)
	return nil
}

func (r *userRepository) UpdateUser(user model.User) error {
	r.users[user.Username] = user
	return nil
}

// NewUserRepository creates a new userRepository instance
func NewUserRepository() UserRepository {
	return &userRepository{
		users: make(map[string]model.User),
	}
}
