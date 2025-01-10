package model

import "time"

type AggregateTotalAmountResponse struct {
	TotalAmount float64 `json:"total_amount"`
	Type        string  `json:"type"`
}

type CountTransactionsPerTypeResponse struct {
	Type             string `json:"type"`
	TransactionCount uint64 `json:"transaction_count"`
}

type AverageTransactionValuePerUserResponse struct {
	UserID                  uint64  `json:"user_id"`
	AverageTransactionValue float64 `json:"average_transaction_value"`
}

type TopUsersByTotalTransactionValueResponse struct {
	UserID      uint64  `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
}

type TransactionDailySummaryResponse struct {
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Purchase uint64    `json:"purchase,omitempty"`
	Refund   uint64    `json:"refund,omitempty"`
}

type GetUsersWithMoreThanXTransactionsResponse struct {
	UserID           uint64 `json:"user_id"`
	TransactionCount uint64 `json:"transaction_count"`
}

type GetUsersWithBothPurchasesAndRefundsResponse struct {
	UserID uint64 `json:"user_id"`
}
