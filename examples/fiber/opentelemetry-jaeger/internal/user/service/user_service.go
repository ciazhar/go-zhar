package service

import (
	"context"
	"errors"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/model"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	userRepository *repository.UserRepository
	tracer         trace.Tracer
}

func NewUserService(
	userRepo *repository.UserRepository,
	tracer trace.Tracer,
) *UserService {
	return &UserService{
		userRepository: userRepo,
		tracer:         tracer,
	}
}

// AddUser adds a new user
func (s *UserService) AddUser(ctx context.Context, user model.User, tracer trace.Span) {
	_, span := s.tracer.Start(ctx, "AddUser")
	defer span.End()
	s.userRepository.AddUser(ctx, user)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	_, span := s.tracer.Start(ctx, "GetUserByUsername")
	defer span.End()
	return s.userRepository.GetUserByUsername(ctx, username)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context, tracer trace.Span) map[string]model.User {
	_, span := s.tracer.Start(trace.ContextWithSpanContext(ctx, tracer.SpanContext()), "GetAllUsersUseCase")
	defer span.End()
	span.RecordError(errors.New("some error"))
	span.AddEvent("some event", trace.WithAttributes(attribute.String("key", "some log")))
	return s.userRepository.GetAllUsers(ctx, span)
}

// DeleteUser deletes a user by username
func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	_, span := s.tracer.Start(ctx, "DeleteUser")
	defer span.End()
	return s.userRepository.DeleteUser(ctx, username)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, user model.User) error {
	_, span := s.tracer.Start(ctx, "UpdateUser")
	defer span.End()
	return s.userRepository.UpdateUser(ctx, user)
}
