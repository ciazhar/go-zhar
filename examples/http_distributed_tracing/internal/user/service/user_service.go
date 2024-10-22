package service

import (
	"context"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/user/model"
	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/user/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	userRepository *repository.UserRepository
	tracer         trace.Tracer
}

func NewUserService(
	userRepo *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepo,
		tracer:         otel.Tracer("UserService"),
	}
}

// AddUser adds a new user
func (s *UserService) AddUser(ctx context.Context, user model.User, parentSpan trace.Span) {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserService_AddUser", trace.WithAttributes(
			attribute.String("username", user.Username),
			attribute.String("email", user.Email),
			attribute.Int("age", user.Age),
		))
	defer span.End()
	s.userRepository.AddUser(ctx, user, span)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string, parentSpan trace.Span) (*model.User, error) {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserService_GetUserByUsername",
		trace.WithAttributes(
			attribute.String("username", username),
		),
	)
	defer span.End()
	return s.userRepository.GetUserByUsername(ctx, username, span)
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context, parentSpan trace.Span) map[string]model.User {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserService_GetAllUsers",
	)
	defer span.End()
	return s.userRepository.GetAllUsers(ctx, span)
}

// DeleteUser deletes a user by username
func (s *UserService) DeleteUser(ctx context.Context, username string, parentSpan trace.Span) error {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserService_DeleteUser",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()
	return s.userRepository.DeleteUser(ctx, username, span)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, user model.User, parentSpan trace.Span) error {
	_, span := s.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserService_UpdateUser",
		trace.WithAttributes(
			attribute.String("username", user.Username),
			attribute.String("email", user.Email),
			attribute.Int("age", user.Age),
		),
	)
	defer span.End()
	return s.userRepository.UpdateUser(ctx, user, span)
}
