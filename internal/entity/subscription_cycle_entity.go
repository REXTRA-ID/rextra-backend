package entity

import (
	"time"

	"github.com/google/uuid"
)

type CycleStatus string

const (
	CycleStatusActive    CycleStatus = "active"
	CycleStatusCompleted CycleStatus = "completed"
	CycleStatusExpired   CycleStatus = "expired"
)

type SubscriptionCycle struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipID uuid.UUID `json:"membership_id" gorm:"type:uuid;not null;index"`
	CycleNumber  int       `json:"cycle_number" gorm:"not null;default:0"`

	PlanName       string `json:"plan_name" gorm:"type:varchar(20);not null"`
	PlanCategory   string `json:"plan_category" gorm:"type:varchar(20);not null"`
	DurationMonths int    `json:"duration_months" gorm:"not null;default:0"`

	PlanDurationID *uuid.UUID `json:"plan_duration_id,omitempty" gorm:"type:uuid;index"`

	StartDate time.Time  `json:"start_date" gorm:"not null"`
	EndDate   *time.Time `json:"end_date,omitempty"`

	AmountPaid     int64  `json:"amount_paid" gorm:"not null;default:0"`
	PaymentChannel string `json:"payment_channel" gorm:"type:varchar(100)"`
	TransactionID  string `json:"transaction_id" gorm:"type:varchar(100);index"`

	Status CycleStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Membership Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	PlanDuration *PlanDuration `json:"plan_duration,omitempty" gorm:"foreignKey:PlanDurationID;references:ID"`
}

func (SubscriptionCycle) TableName() string {
	return "subscription_cycles"
}

func NewSubscriptionCycle(
	membershipID uuid.UUID,
	planDurationID *uuid.UUID,
	planName, planCategory string,
	durationMonths int,
	startDate time.Time,
	amountPaid int64,
	paymentChannel, transactionID string,
	cycleNumber int,
) SubscriptionCycle {
	return SubscriptionCycle{
		MembershipID:   membershipID,
		PlanDurationID: planDurationID,
		PlanName:       planName,
		PlanCategory:   planCategory,
		DurationMonths: durationMonths,
		StartDate:      startDate,
		AmountPaid:     amountPaid,
		PaymentChannel: paymentChannel,
		TransactionID:  transactionID,
		CycleNumber:    cycleNumber,
		Status:         CycleStatusActive,
		CreatedAt:      time.Now().UTC(),
	}
}
