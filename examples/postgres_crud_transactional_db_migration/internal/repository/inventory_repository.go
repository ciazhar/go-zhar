package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepository interface {
	Update(ctx context.Context, tx pgx.Tx, quantity, productID int) error
}

type PgxInventoryRepository struct {
	pool *pgxpool.Pool
}

func NewPgxInventoryRepository(pool *pgxpool.Pool) *PgxInventoryRepository {
	return &PgxInventoryRepository{
		pool: pool,
	}
}

func (r *PgxInventoryRepository) Update(ctx context.Context, tx pgx.Tx, quantity, productID int) error {
	_, err := tx.Exec(ctx, "UPDATE inventory SET stock = stock - $1 WHERE product_id = $2 AND stock >= $1", quantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update inventory for product %d: %w", productID, err)
	}
	return nil
}
