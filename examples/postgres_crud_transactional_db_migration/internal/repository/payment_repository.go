package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	Create(ctx context.Context, tx pgx.Tx, orderID int, paymentAmount int, status string) error
	Delete(ctx context.Context, tx pgx.Tx, paymentID int) error
}

type PgxPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPgxPaymentRepository(pool *pgxpool.Pool) *PgxPaymentRepository {
	return &PgxPaymentRepository{
		pool: pool,
	}
}

func (r *PgxPaymentRepository) Create(ctx context.Context, tx pgx.Tx, orderID int, paymentAmount int, status string) error {
	_, err := tx.Exec(ctx, "INSERT INTO payments (order_id, amount, status) VALUES ($1, $2, $3)", orderID, paymentAmount, status)
	if err != nil {
		return fmt.Errorf("failed to process payment: %w", err)
	}
	return nil
}

func (r *PgxPaymentRepository) Delete(ctx context.Context, tx pgx.Tx, paymentID int) error {
	_, err := tx.Exec(ctx, "DELETE FROM payments WHERE id = $1", paymentID)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	return nil
}
