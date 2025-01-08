package repository

import (
	"context"
	"errors"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]model.Product, error)
	CreateProduct(ctx context.Context, name string, price float64, stock int) (int, error)
}

type PgxProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *PgxProductRepository {
	return &PgxProductRepository{
		pool: pool,
	}
}

func (r *PgxProductRepository) GetProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, price, stock FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PgxProductRepository) CreateProduct(ctx context.Context, name string, price float64, stock int) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`, name, price, stock).Scan(&id)
	return id, err
}

func (r *PgxProductRepository) GetProductByID(ctx context.Context, tx pgx.Tx, productID int) (*model.Product, error) {
	query := `SELECT id, name, stock, price FROM products WHERE id = $1`
	product := &model.Product{}
	err := tx.QueryRow(ctx, query, productID).Scan(&product.ID, &product.Name, &product.Stock, &product.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *PgxProductRepository) AdjustStock(ctx context.Context, tx pgx.Tx, productID int, quantity int) error {
	query := `UPDATE products SET stock = stock + $1 WHERE id = $2`
	_, err := tx.Exec(ctx, query, quantity, productID)
	return err
}
