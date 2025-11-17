package entity

import (
	"time"

	"github.com/google/uuid"
)

type Memberships struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID           uuid.UUID    `json:"user_id" gorm:"not null"`
	MembershipStatus EnumPlanName `json:"membership_status" gorm:"type:varchar(20)"`

	Plan   MembershipPlans `gorm:"foreignKey:ID;references:PlanID"`
	PlanID uuid.UUID       `json:"plan_id"`

	Duration   MembershipDuration `gorm:"foreignKey:ID;references:DurationID"`
	DurationID uuid.UUID          `json:"duration_id"`

	CurrentTokenBalance int       `json:"current_token_balance" gorm:"default:0"`
	CurrentPoinBalance  int       `json:"current_poin_balance" gorm:"default:0"`
	StartedAt           time.Time `json:"started_at"`
	ExpiredAt           time.Time `json:"expired_at"`
	IsActive            bool      `json:"is_active"`
	AutoRenew           bool      `json:"auto_renew"`

	Timestamp
}

func (m *Memberships) TableName() string {
	return "memberships"
}

func NewMembership(userId uuid.UUID, plan *MembershipPlans, duration *MembershipDuration) Memberships {
	return Memberships{
		UserID:           userId,
		MembershipStatus: PLANSTARTER,
		PlanID:           plan.ID,
		DurationID:       duration.ID,
		StartedAt:        time.Now().UTC(),
		ExpiredAt:        time.Now().Add(time.Duration(time.Now().Day() * 30)).UTC(),
	}
}

func (m *Memberships) IsUserMembershipExpired() bool {
	return time.Now().After(m.ExpiredAt)
}
