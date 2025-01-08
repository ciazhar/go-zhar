package model

import "time"

// Customer represents a customer entity.
type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Order represents an order entity.
type Order struct {
	ID          int       `json:"id"`
	CustomerID  int       `json:"customer_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderItem represents an order item entity.
type OrderItem struct {
	ID         int       `json:"id"`
	OrderID    int       `json:"order_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

// Product represents a product entity.
type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Payment represents a payment entity.
type Payment struct {
	ID            int       `json:"id"`
	OrderID       int       `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Shipment represents a shipment entity.
type Shipment struct {
	ID             int       `json:"id"`
	OrderID        int       `json:"order_id"`
	TrackingNumber string    `json:"tracking_number"`
	Carrier        string    `json:"carrier"`
	Status         string    `json:"status"`
	ShippedAt      time.Time `json:"shipped_at"`
	DeliveredAt    time.Time `json:"delivered_at"`
}
