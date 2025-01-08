package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
	CreateOrder(ctx context.Context, customerID int, status string) (int, error)
	AddOrderItem(ctx context.Context, orderID, productID, quantity int, price, totalPrice float64) error
	UpdateOrderTotal(ctx context.Context, orderID int, totalAmount float64) error
}

type PgxOrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *PgxOrderRepository {
	return &PgxOrderRepository{
		pool: pool,
	}
}

func (r *PgxOrderRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (r *PgxOrderRepository) CreateOrder(ctx context.Context, tx pgx.Tx, customerID int, status string) (int, error) {
	var id int
	err := tx.QueryRow(ctx, `INSERT INTO orders (customer_id, status) VALUES ($1, $2) RETURNING id`, customerID, status).Scan(&id)
	return id, err
}

func (r *PgxOrderRepository) AddOrderItem(ctx context.Context, tx pgx.Tx, orderID, productID, quantity int, price, totalPrice float64) error {
	_, err := tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, price, total_price) VALUES ($1, $2, $3, $4, $5)`, orderID, productID, quantity, price, totalPrice)
	return err
}

func (r *PgxOrderRepository) UpdateOrderTotal(ctx context.Context, tx pgx.Tx, orderID int, totalAmount float64) error {
	_, err := tx.Exec(ctx, `UPDATE orders SET total_amount = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, totalAmount, orderID)
	return err
}

func (r *PgxOrderRepository) UpdateOrderStatus(ctx context.Context, tx pgx.Tx, orderID int, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := tx.Exec(ctx, query, status, orderID)
	return err
}
