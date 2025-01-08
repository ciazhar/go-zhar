package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, name, email string) (int, error)
}

type PgxCustomerRepository struct {
	pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) *PgxCustomerRepository {
	return &PgxCustomerRepository{
		pool: pool,
	}
}

func (r *PgxCustomerRepository) CreateCustomer(ctx context.Context, name, email string) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, `INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING id`, name, email).Scan(&id)
	return id, err
}
