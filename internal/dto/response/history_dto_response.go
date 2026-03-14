package dto_response

import "time"

type TransactionType string

const (
	TransactionTypeClub  TransactionType = "club"
	TransactionTypeToken TransactionType = "token"
)

type HistoryTransactionItem struct {
	TransactionType TransactionType `json:"transaction_type"`
	TransactionID   string          `json:"transaction_id"`
	Status          string          `json:"status"`
	StatusLabel     string          `json:"status_label"`
	Amount          int64           `json:"amount"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
	PaidAt          *time.Time      `json:"paid_at,omitempty"`
	PaymentMethod   *string         `json:"payment_method,omitempty"`
}

type HistoryTransactionListResponse struct {
	Items    []HistoryTransactionItem `json:"items"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	HasMore  bool                     `json:"has_more"`
}
