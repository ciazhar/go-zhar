package user

import (
	"context"
	"go.opentelemetry.io/otel"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) GetUsersWithPagination(ctx context.Context, page, limit int) ([]response.User, int64, error) {
	var (
		reqCtx, span = otel.Tracer("repository").Start(ctx, "UserRepository.GetUsersWithPagination")
		deferFn      = func() { span.End() }
		log          = logger.FromContext(reqCtx).With().Int("page", page).Int("limit", limit).Logger()
		offset       = (page - 1) * limit
		resp         = make([]response.User, 0)
	)
	defer deferFn()

	rows, err := r.pg.Query(reqCtx, queryGetUsersWithPagination, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to get users")
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u response.User
		if err = rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("failed to scan user")
			return nil, 0, err
		}
		resp = append(resp, u)
	}

	var total int64
	err = r.pg.QueryRow(reqCtx, queryCountUsers).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("failed to count users")
		return nil, 0, err
	}

	return resp, total, nil
}
