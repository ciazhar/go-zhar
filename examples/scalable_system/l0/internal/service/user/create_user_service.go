package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (u userService) CreateUser(ctx context.Context, req request.CreateUserBodyRequest) error {
	var (
		reqCtx, span = otel.Tracer("service").Start(ctx, "UserService.CreateUser")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Any("req", req).Logger()
	)
	defer deferFn()

	err := u.repo.CreateUser(reqCtx, req)
	if err != nil {
		log.Err(err).Msg("failed to insert user to DB")
		return err
	}

	return nil
}
