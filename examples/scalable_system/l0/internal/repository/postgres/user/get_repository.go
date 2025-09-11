package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/response"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r UserRepository) GetUsers(ctx context.Context, page, limit int) ([]response.User, int64, error) {
	var (
		log    = logger.FromContext(ctx).With().Int("page", page).Int("limit", limit).Logger()
		offset = (page - 1) * limit
		resp   = make([]response.User, 0)
	)

	rows, err := r.pg.Query(ctx, queryGetUsersWithPagination, limit, offset)
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
	err = r.pg.QueryRow(ctx, queryCountUsers).Scan(&total)
	if err != nil {
		log.Error().Err(err).Msg("failed to count users")
		return nil, 0, err
	}

	return resp, total, nil
}
