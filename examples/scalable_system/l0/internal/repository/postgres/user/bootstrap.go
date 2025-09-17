package user

import (
	"github.com/ciazhar/go-zhar/pkg/postgres"
)

type UserRepository struct {
	pg postgres.PgxPool
}

func NewUserRepository(pg postgres.PgxPool) *UserRepository {
	return &UserRepository{
		pg: pg,
	}
}
