package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShipmentRepository interface {
	CreateShipment(ctx context.Context, orderID int, trackingNumber, carrier, status string) error
	UpdateShipmentStatus(ctx context.Context, orderID int, status string) error
}

type PgxShipmentRepository struct {
	pool *pgxpool.Pool
}

func NewShipmentRepository(pool *pgxpool.Pool) *PgxShipmentRepository {
	return &PgxShipmentRepository{
		pool: pool,
	}
}

func (r *PgxShipmentRepository) CreateShipment(ctx context.Context, orderID int, trackingNumber, carrier, status string) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO shipments (order_id, tracking_number, carrier, status, shipped_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)`, orderID, trackingNumber, carrier, status)
	return err
}

func (r *PgxShipmentRepository) UpdateShipmentStatus(ctx context.Context, orderID int, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE shipments SET status = $1, delivered_at = CURRENT_TIMESTAMP WHERE order_id = $2`, status, orderID)
	return err
}
