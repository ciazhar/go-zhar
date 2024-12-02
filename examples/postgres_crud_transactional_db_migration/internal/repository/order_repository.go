package repository

import (
	"context"
	"fmt"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, tx pgx.Tx, customerName string) (int, error)
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, orderID int) (*model.Order, error)
	DeleteOrderByID(ctx context.Context, orderID int) error
}

type PgxOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPgxOrderRepository(pool *pgxpool.Pool) *PgxOrderRepository {
	return &PgxOrderRepository{
		pool: pool,
	}
}

func (r *PgxOrderRepository) Create(ctx context.Context, tx pgx.Tx, customerName string) (int, error) {
	var orderID int
	err := tx.QueryRow(ctx, "INSERT INTO orders (customer_name) VALUES ($1) RETURNING id", customerName).Scan(&orderID)
	if err != nil {
		return orderID, fmt.Errorf("failed to create order: %w", err)
	}
	return 0, nil
}

func (r *PgxOrderRepository) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}

func (r *PgxOrderRepository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, customer_name, created_at FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PgxOrderRepository) GetOrderByID(ctx context.Context, orderID int) (*model.Order, error) {
	var order model.Order
	err := r.pool.QueryRow(ctx, "SELECT id, customer_name, created_at FROM orders WHERE id = $1", orderID).Scan(&order.ID, &order.CustomerName, &order.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	return &order, nil
}

func (r *PgxOrderRepository) DeleteOrderByID(ctx context.Context, orderID int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
