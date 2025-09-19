package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar-scalable-system-l0/internal/model/response"

	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) GetUserByID(ctx context.Context, id string) (*response.User, error) {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.GetUserByID")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Str("id", id).Logger()
		resp         = new(response.User)
	)
	defer deferFn()

	err := r.pg.QueryRow(reqCtx, queryGetUserByID, id).
		Scan(&resp.ID, &resp.Username, &resp.Email, &resp.FullName, &resp.CreatedAt, &resp.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by ID")
		return nil, err
	}

	return resp, nil
}
