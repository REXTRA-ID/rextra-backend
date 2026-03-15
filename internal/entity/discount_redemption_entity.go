package entity

import (
	"time"

	"github.com/google/uuid"
)

type RedemptionStatus string

const (
	RedemptionStatusApplied  RedemptionStatus = "APPLIED"
	RedemptionStatusReversed RedemptionStatus = "REVERSED"
)

type DiscountRedemption struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DiscountID uuid.UUID `json:"discount_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	TransactionID string `json:"transaction_id" gorm:"type:varchar(100);not null;index"`

	CodeSnapshot  string  `json:"code_snapshot" gorm:"type:varchar(50);not null"`
	PlanSnapshot  *string `json:"plan_snapshot,omitempty" gorm:"type:varchar(20)"`
	UserName      string  `json:"user_name" gorm:"type:varchar(200)"`
	AppliesToType string  `json:"applies_to_type" gorm:"type:varchar(20)"`

	SubtotalAmount int64 `json:"subtotal_amount" gorm:"not null;default:0"`
	DiscountAmount int64 `json:"discount_amount" gorm:"not null;default:0"`
	FinalAmount    int64 `json:"final_amount" gorm:"not null;default:0"`

	Status    RedemptionStatus `json:"status" gorm:"type:varchar(20);not null;default:'APPLIED'"`
	AppliedAt time.Time        `json:"applied_at" gorm:"not null"`

	ReversedAt    *time.Time `json:"reversed_at,omitempty"`
	ReverseReason *string    `json:"reverse_reason,omitempty" gorm:"type:text"`

	Discount Discounts `json:"discount,omitempty" gorm:"foreignKey:DiscountID;references:ID"`
}

func (DiscountRedemption) TableName() string {
	return "discount_redemptions"
}

func NewDiscountRedemption(
	discountID uuid.UUID,
	userID uuid.UUID,
	transactionID string,
	codeSnapshot string,
	planSnapshot *string,
	userName string,
	appliesToType string,
	subtotalAmount, discountAmount, finalAmount int64,
) DiscountRedemption {
	return DiscountRedemption{
		DiscountID:     discountID,
		UserID:         userID,
		TransactionID:  transactionID,
		CodeSnapshot:   codeSnapshot,
		PlanSnapshot:   planSnapshot,
		UserName:       userName,
		AppliesToType:  appliesToType,
		SubtotalAmount: subtotalAmount,
		DiscountAmount: discountAmount,
		FinalAmount:    finalAmount,
		Status:         RedemptionStatusApplied,
		AppliedAt:      time.Now().UTC(),
	}
}

func (r *DiscountRedemption) Reverse(reason string) {
	now := time.Now().UTC()
	r.Status = RedemptionStatusReversed
	r.ReversedAt = &now
	r.ReverseReason = &reason
}
