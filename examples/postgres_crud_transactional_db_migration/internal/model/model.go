package model

import "time"

type OrderRequest struct {
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Amount       float64     `json:"amount"`
}
type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type Inventory struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

type Payment struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
