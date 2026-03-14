package entity

import (
	"time"

	"github.com/google/uuid"
)

type PoinTransactionType string
type PoinSource string

const (
	PoinTypeEarn       PoinTransactionType = "earn"
	PoinTypeSpend      PoinTransactionType = "spend"
	PoinTypeAdjustment PoinTransactionType = "adjustment"

	PoinSourceMembershipPurchase PoinSource = "membership_purchase"
	PoinSourceMissionComplete    PoinSource = "mission_complete"
	PoinSourceEventReward        PoinSource = "event_reward"
	PoinSourceRedeemReward       PoinSource = "redeem_reward"
	PoinSourceAdminAdjustment    PoinSource = "admin_adjustment"
)

type PoinTransactions struct {
	ID                uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID            uuid.UUID           `json:"user_id" gorm:"type:uuid;not null;index"`
	MembershipID      uuid.UUID           `json:"membership_id" gorm:"type:uuid;not null;index"`
	TransactionType   PoinTransactionType `json:"transaction_type" gorm:"type:varchar(20);not null"`
	PoinAmount        int                 `json:"poin_amount" gorm:"not null"`
	PoinBalanceBefore int                 `json:"poin_balance_before" gorm:"not null"`
	PoinBalanceAfter  int                 `json:"poin_balance_after" gorm:"not null"`
	Source            PoinSource          `json:"source" gorm:"type:varchar(50);not null"`
	ReferenceID       *uuid.UUID          `json:"reference_id,omitempty" gorm:"type:uuid"`
	Description       string              `json:"description" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Membership Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (PoinTransactions) TableName() string {
	return "poin_transactions"
}

func NewPoinTransaction(
	userID uuid.UUID,
	membershipID uuid.UUID,
	transactionType PoinTransactionType,
	source PoinSource,
	currentBalance int,
	poinAmount int,
	referenceID *uuid.UUID,
	description string,
) PoinTransactions {
	var poinBalanceAfter int

	switch transactionType {
	case PoinTypeEarn, PoinTypeAdjustment:
		poinBalanceAfter = currentBalance + poinAmount
	case PoinTypeSpend:
		poinBalanceAfter = currentBalance - poinAmount
	default:
		poinBalanceAfter = currentBalance
	}

	return PoinTransactions{
		UserID:            userID,
		MembershipID:      membershipID,
		TransactionType:   transactionType,
		PoinAmount:        poinAmount,
		PoinBalanceBefore: currentBalance,
		PoinBalanceAfter:  poinBalanceAfter,
		Source:            source,
		ReferenceID:       referenceID,
		Description:       description,
		CreatedAt:         time.Now().UTC(),
	}
}
