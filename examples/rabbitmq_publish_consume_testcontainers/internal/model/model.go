package model

type OrderStatusRequest struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

type PaymentReminderRequest struct {
	OrderId  string `json:"order_id"`
	Reminder string `json:"reminder"`
}
