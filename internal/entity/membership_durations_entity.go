package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipDuration struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DurationMonth        int       `json:"duration_months"`
	TokenBonusPercentage float64   `json:"token_bonus_percentage" gorm:"default:0"`
	RextraPoinMultiplier int       `json:"rextra_poin_multiplier" gorm:"default:1"`
	IsActive             bool      `json:"is_active" gorm:"type:boolean;default:true"`

	Timestamp
}

func (m *MembershipDuration) TableName() string {
	return "membership_duration"
}

func NewMembershipDuration(durationMonth int, tokenBonusPercentage float64, rextraPoinMultiplier int) MembershipDuration {
	return MembershipDuration{
		DurationMonth:        durationMonth,
		TokenBonusPercentage: tokenBonusPercentage,
		RextraPoinMultiplier: rextraPoinMultiplier,
	}
}

func (m *MembershipDuration) BeforeCreate(tx *gorm.DB) (err error) {
	time := time.Now().UTC()
	m.CreatedAt = time
	m.UpdatedAt = time
	return nil
}
