package entity

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionCycle struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID `json:"user_id" gorm:"not null"`
	MembershipID   uuid.UUID `json:"membership_id" gorm:"not null"`
	PlanName       string    `json:"plan_name"`
	DurationMonths int       `json:"duration_months"`
	AmountPaid     float64   `json:"amount_paid"`
	PaymentChannel string    `json:"payment_channel"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Timestamp
}

func (s *SubscriptionCycle) TableName() string {
	return "subscription_cycles"
}

func NewSubscriptionCycle(userId, membershipId uuid.UUID, planName string, durationMonth int, amountPaid float64, paymentChannel string, startDate, endDate time.Time) SubscriptionCycle {
	return SubscriptionCycle{
		UserID:         userId,
		MembershipID:   membershipId,
		PlanName:       planName,
		DurationMonths: durationMonth,
		AmountPaid:     amountPaid,
		PaymentChannel: paymentChannel,
		StartDate:      startDate,
		EndDate:        endDate,
	}
}
