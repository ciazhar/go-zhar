package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/internal/user/model"
	model2 "github.com/ciazhar/go-zhar/examples/fiber/opentelemetry-jaeger/pkg/model"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

type UserHTTPRepository struct {
	client         *http.Client
	tracer         trace.Tracer
	userServiceURL string
}

func NewUserHTTPRepository(client *http.Client, tracer trace.Tracer) *UserHTTPRepository {
	return &UserHTTPRepository{
		client:         client,
		tracer:         tracer,
		userServiceURL: viper.GetString("application.user-service-url"),
	}
}

func (r *UserHTTPRepository) GetUserByUsername(ctx context.Context, parentSpan trace.Span, username string) (model.User, error) {
	var user model.User

	ctx, span := r.tracer.Start(
		trace.ContextWithSpanContext(ctx, parentSpan.SpanContext()),
		"UserHTTPRepository_GetUserByUsername",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", r.userServiceURL, username), nil)
	if err != nil {
		return user, fmt.Errorf("failed to create request: %w", err)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := r.client.Do(req)
	if err != nil {
		return user, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("failed to fetch user: status code %d", resp.StatusCode)
	}

	var userResponse model2.Response
	if err := decodeBody(resp.Body, &userResponse); err != nil {
		return user, fmt.Errorf("failed to decode response body: %w", err)
	}

	userData, err := json.Marshal(userResponse.Data)
	if err != nil {
		return user, fmt.Errorf("failed to marshal user data: %w", err)
	}

	if err := json.Unmarshal(userData, &user); err != nil {
		return user, fmt.Errorf("invalid response data format: %w", err)
	}

	return user, nil
}

func decodeBody(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(v)
}
