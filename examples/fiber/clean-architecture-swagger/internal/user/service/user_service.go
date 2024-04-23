package service

import (
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/clean-architecture-swagger/internal/user/repository"
)

// UserService interface defines methods for user service
type UserService interface {
	AddUser(user model.User)
	GetUserByUsername(username string) (*model.User, error)
	GetAllUsers() map[string]model.User
	DeleteUser(username string) error
	UpdateUser(user model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepository: userRepo}
}

// AddUser adds a new user
func (s *userService) AddUser(user model.User) {
	s.userRepository.AddUser(user)
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() map[string]model.User {
	return s.userRepository.GetAllUsers()
}

// DeleteUser deletes a user by username
func (s *userService) DeleteUser(username string) error {
	return s.userRepository.DeleteUser(username)
}

// UpdateUser updates a user
func (s *userService) UpdateUser(user model.User) error {
	return s.userRepository.UpdateUser(user)
}
