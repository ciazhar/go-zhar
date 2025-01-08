package service

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/model"
	"github.com/ciazhar/go-start-small/examples/postgres_crud_transactional_db_migration/internal/repository"
)

type ProductService struct {
	Repo *repository.PgxProductRepository
}

func NewProductService(repo *repository.PgxProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]model.Product, error) {
	return s.Repo.GetProducts(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, price float64, stock int) (int, error) {
	return s.Repo.CreateProduct(ctx, name, price, stock)
}
