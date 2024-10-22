package repository

import (
	"context"
	"fmt"

	"github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/user/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UserRepository represents a repository for managing users
type UserRepository struct {
	users  map[string]model.User
	tracer trace.Tracer
}

// AddUser adds a new user to the repository
func (r *UserRepository) AddUser(ctx context.Context, user model.User, parentSpan trace.Span) {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserRepository_AddUser", trace.WithAttributes(
			attribute.String("username", user.Username),
			attribute.String("email", user.Email),
			attribute.Int("age", user.Age),
		))
	defer span.End()
	r.users[user.Username] = user
}

// GetUserByUsername retrieves a user by their username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string, parentSpan trace.Span) (*model.User, error) {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserRepository_GetUserByUsername",
		trace.WithAttributes(
			attribute.String("username", username),
		),
	)
	defer span.End()
	for _, user := range r.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// GetAllUsers retrieves all users from the repository
func (r *UserRepository) GetAllUsers(ctx context.Context, parentSpan trace.Span) map[string]model.User {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserRepository_GetAllUsers",
	)
	defer span.End()
	return r.users
}

// DeleteUser deletes a user from the repository
func (r *UserRepository) DeleteUser(ctx context.Context, username string, parentSpan trace.Span) error {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserRepository_DeleteUser",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()
	_, ok := r.users[username]
	if !ok {
		return fmt.Errorf("user not found for username: %s", username)
	}

	delete(r.users, username)
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user model.User, parentSpan trace.Span) error {
	_, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserRepository_UpdateUser",
		trace.WithAttributes(
			attribute.String("username", user.Username),
			attribute.String("email", user.Email),
			attribute.Int("age", user.Age),
		),
	)
	defer span.End()
	_, ok := r.users[user.Username]
	if !ok {
		return fmt.Errorf("user not found for username: %s", user.Username)
	}
	r.users[user.Username] = user
	return nil
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[string]model.User),
		tracer: otel.Tracer("UserRepository"),
	}
}
