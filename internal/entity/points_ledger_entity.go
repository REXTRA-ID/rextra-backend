package entity

import (
	"time"

	"github.com/google/uuid"
)

type PointsMutationType string
type PointsLedgerSource string

const (
	PointsMutationEarn  PointsMutationType = "earn"
	PointsMutationSpend PointsMutationType = "spend"

	PointsSourceMembership PointsLedgerSource = "membership"
	PointsSourceRedemption PointsLedgerSource = "redemption"
	PointsSourceActivity   PointsLedgerSource = "activity"
	PointsSourceAdjustment PointsLedgerSource = "adjustment"
)

type PointsLedger struct {
	ID           uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipID uuid.UUID          `json:"membership_id" gorm:"type:uuid;not null;index"`
	Amount       int                `json:"amount" gorm:"not null"`
	MutationType PointsMutationType `json:"mutation_type" gorm:"type:varchar(20);not null"`
	Source       PointsLedgerSource `json:"source" gorm:"type:varchar(30);not null"`
	ReferenceID  *string            `json:"reference_id,omitempty" gorm:"type:varchar(100)"`
	Description  string             `json:"description" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Membership Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (PointsLedger) TableName() string {
	return "points_ledger"
}

func NewPointsLedgerEntry(
	membershipID uuid.UUID,
	amount int,
	mutationType PointsMutationType,
	source PointsLedgerSource,
	referenceID *string,
	description string,
) PointsLedger {
	return PointsLedger{
		MembershipID: membershipID,
		Amount:       amount,
		MutationType: mutationType,
		Source:       source,
		ReferenceID:  referenceID,
		Description:  description,
		CreatedAt:    time.Now().UTC(),
	}
}
