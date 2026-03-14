package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	TransactionType string

	Source string
)

const (
	EARN       TransactionType = "earn"
	REDEEM     TransactionType = "redeem"
	ADJUSTMENT TransactionType = "adjustment"

	MEMBERSHIPPURCHASE Source = "membership_purchase"
	MISSIONCOMPLETE    Source = "mission_complete"
	EVENTREWARD        Source = "event_reward"
	REDEEMREWARD       Source = "redeem_reward"
)

type PoinTransactions struct {
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID            uuid.UUID       `json:"user_id"`
	TransactionType   TransactionType `json:"transaction_type"`
	PointAmount       int             `json:"poin_amount"`
	PoinBalanceBefore int             `json:"poin_balance_before"`
	PointBalanceAfter int             `json:"point_balance_after"`
	Source            Source          `json:"source"`
	ReferenceID       uuid.UUID       `json:"reference_id"` // perlu ditanyakan
	Description       string          `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}

func (p *PoinTransactions) TableName() string {
	return "poin_transactions"
}

func NewPoinTransaction(userId uuid.UUID, transactionType string, source string, currentPoin, poinAmount int) PoinTransactions {
	var poinAfter int

	if transactionType == string(EARN) {
		poinAfter += poinAfter
	} else if transactionType == string(REDEEM) {
		poinAfter -= poinAfter
	}

	return PoinTransactions{
		UserID:            userId,
		TransactionType:   TransactionType(transactionType),
		PointAmount:       poinAmount,
		PoinBalanceBefore: currentPoin,
		PointBalanceAfter: poinAfter,
		Source:            Source(source),
		Description:       "",
		CreatedAt:         time.Now().UTC(),
	}
}
