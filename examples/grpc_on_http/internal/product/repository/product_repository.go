package repository

import (
	"context"
	"errors"
	"github.com/ciazhar/go-start-small/examples/grpc_on_http/internal/product/model"
)

// ProductRepository defines the interface
type ProductRepository interface {
	GetProductByID(ctx context.Context, id int32) (*model.Product, error)
}

// DummyProductRepository is a dummy implementation of ProductRepository
type DummyProductRepository struct {
	data map[int32]*model.Product
}

// NewDummyProductRepository creates a new instance of DummyProductRepository
func NewDummyProductRepository() ProductRepository {
	return &DummyProductRepository{
		data: map[int32]*model.Product{
			1: {ID: 1, Name: "Kopi Susu", Price: 18000, Image: "https://example.com/image/kopi.jpg"},
			2: {ID: 2, Name: "Teh Manis", Price: 12000, Image: "https://example.com/image/teh.jpg"},
			3: {ID: 3, Name: "Air Mineral", Price: 8000, Image: "https://example.com/image/air.jpg"},
		},
	}
}

// GetProductByID fetches product by its ID from dummy data
func (r *DummyProductRepository) GetProductByID(ctx context.Context, id int32) (*model.Product, error) {
	product, ok := r.data[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}
