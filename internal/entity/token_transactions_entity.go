package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	TokenTransactionType string

	ReferenceType string
)

const (
	TOKENPURCHASEMEMBERSHIP TokenTransactionType = "purchase_membership"
	TOKENPURCHASESTANDALONE TokenTransactionType = "purchase_standalone"
	TOKENREFUND             TokenTransactionType = "refund"
	TOKENMONTHLYREFILL      TokenTransactionType = "monthly_refill"

	REFPAYMENT            ReferenceType = "payment"
	REFFEATURE_USAHE      ReferenceType = "feature_usage"
	REFMEMBERSHIP_RENEWAL ReferenceType = "membership_renewal"
	REFMONTHLYREFILL      ReferenceType = "monthly_refill"
)

type TokenTransaction struct {
	ID                 uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID             uuid.UUID            `json:"user_id"`
	TransactionType    TokenTransactionType `json:"transaction_type" gorm:"type:varchar(30)"`
	TokenAmount        int                  `json:"token_amount"`
	TokenBalanceBefore int                  `json:"token_balance_before"`
	TokenBalanceAfter  int                  `json:"token_balance_after"`
	ReferenceType      ReferenceType        `json:"reference_type"`
	ReferenceID        uuid.UUID            `json:"reference_id"` //masih ambigu
	Description        string               `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}

func (t *TokenTransaction) TableName() string {
	return "token_transactions"
}
