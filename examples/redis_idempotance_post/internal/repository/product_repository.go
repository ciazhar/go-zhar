package repository

import (
	"errors"
	"github.com/ciazhar/go-start-small/examples/redis_idempotance_post/internal/model"
)

var productDB = make(map[string]model.Product)

var DB = struct {
	FindProductBySKU func(string) (*model.Product, error)
	InsertProduct    func(request model.CreateProductRequest) error
}{
	FindProductBySKU: func(sku string) (*model.Product, error) {
		if p, ok := productDB[sku]; ok {
			return &p, nil
		}
		return nil, errors.New("not found")
	},
	InsertProduct: func(req model.CreateProductRequest) error {
		if _, exists := productDB[req.SKU]; exists {
			return errors.New("already exists")
		}
		productDB[req.SKU] = model.Product{Name: req.Name, SKU: req.SKU}
		return nil
	},
}
