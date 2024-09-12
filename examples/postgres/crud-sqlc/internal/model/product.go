package model

// CreateProductRequest represents the request body for creating a new product
type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// UpdateProductPriceRequest represents the request body for updating a product's price
type UpdateProductPriceRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Product represents a product
type Product struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt int64   `json:"created_at"`
}
