package model

import "github.com/ciazhar/go-start-small/examples/http_distributed_tracing/internal/user/model"

type Order struct {
	OrderID   string `json:"order_id"`
	OrderDate string `json:"order_date"`
	Username  string `json:"username"`
}

type OrderExtended struct {
	OrderID   string     `json:"order_id"`
	OrderDate string     `json:"order_date"`
	User      model.User `json:"user"`
}
