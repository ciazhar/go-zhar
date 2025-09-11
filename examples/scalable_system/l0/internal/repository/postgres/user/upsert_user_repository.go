package user

import (
	"context"

	"github.com/ciazhar/go-zhar/examples/scalable_system/l0/internal/model/request"
	"github.com/ciazhar/go-zhar/pkg/logger"
)

func (r *UserRepository) UpsertUserByID(ctx context.Context, req request.UpsertUserBodyRequest) error {
	var (
		log = logger.FromContext(ctx).With().Any("req", req).Logger()
	)

	if _, err := r.pg.Exec(ctx, queryUpsertUser, req.Id, req.Username, req.Email, req.Password, req.FullName); err != nil {
		log.Error().Err(err).Msg("failed to upsert user")
		return err
	}

	return nil
}
