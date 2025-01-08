package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, orderID int, method string, amount float64, status string) error
	GetProductByID(ctx context.Context, productID int)
}

type PgxPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PgxPaymentRepository {
	return &PgxPaymentRepository{
		pool: pool,
	}
}

func (r *PgxPaymentRepository) CreatePayment(ctx context.Context, orderID int, method string, amount float64, status string) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO payments (order_id, payment_method, amount, status) VALUES ($1, $2, $3, $4)`, orderID, method, amount, status)
	return err
}

func (r *PgxPaymentRepository) ProcessPayment(ctx context.Context, tx pgx.Tx, orderID int, method string, amount float64, status string) error {
	query := `INSERT INTO payments (order_id, payment_method, amount, status) VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(ctx, query, orderID, method, amount, status)
	return err
}
