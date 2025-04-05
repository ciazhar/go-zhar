package model

type CreateProductRequest struct {
	Name string `json:"name"`
	SKU  string `json:"sku"` // SKU sebagai kolom unik
}

type Product struct {
	Name string
	SKU  string
}
