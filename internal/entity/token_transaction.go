package entity

import "github.com/google/uuid"

type TransactionTokenType string

const (
	REDEEM      TransactionTokenType = "REDEEM"
	FEATURE_USE TransactionTokenType = "FEATURE_USE"
)

type TokenTransaction struct {
	ID     uuid.UUID            `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	UserID uuid.UUID            `json:"user_id" gorm:"not null"`
	Amount int                  `json:"amount" gorm:"not null"`
	Type   TransactionTokenType `json:"type" gorm:"not null"`
}

func (t *TokenTransaction) TableName() string {
	return "token_transactions"
}
