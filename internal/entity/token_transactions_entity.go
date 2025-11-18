package entity

import (
	"time"

	"github.com/google/uuid"
)

type TokenTransactionType string

const (
	TOKENPURCHASEMEMBERSHIP TokenTransactionType = "purchase_membership"
	TOKENPURCHASESTANDALONE TokenTransactionType = "purchase_standalone"
	TOKENUSAGE              TokenTransactionType = "usage"
	TOKENREFUND             TokenTransactionType = "refund"
	TOKENMONTHLYREFILL      TokenTransactionType = "monthly_refill"
)

type TokenTransaction struct {
	ID                 uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID             uuid.UUID            `json:"user_id"`
	TransactionType    TokenTransactionType `json:"transaction_type" gorm:"type:varchar(30)"`
	TokenAmount        int                  `json:"token_amount"`
	TokenBalanceBefore int                  `json:"token_balance_before"`
	TokenBalanceAfter  int                  `json:"token_balance_after"`
	ReferenceType      string               `json:"reference_type" gorm:"type:varchar(50)"`
	ReferenceID        uuid.UUID            `json:"reference_id"` //masih ambigu
	Description        string               `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}

func (t *TokenTransaction) TableName() string {
	return "token_transactions"
}

func NewTokenTransaction(userId uuid.UUID, userMembership *Memberships, transactionType string, tokenAmount int, description string) TokenTransaction {
	currentToken := userMembership.CurrentTokenBalance

	if transactionType == string(TOKENUSAGE) {
		userMembership.UseMembershipToken(tokenAmount)
	} else {
		userMembership.AddMembershipToken(tokenAmount)
	}

	return TokenTransaction{
		UserID:             userId,
		TransactionType:    TokenTransactionType(transactionType),
		TokenAmount:        tokenAmount,
		TokenBalanceBefore: currentToken,
		TokenBalanceAfter:  userMembership.CurrentTokenBalance,
		Description:        description,
		CreatedAt:          time.Now().UTC(),
	}
}
