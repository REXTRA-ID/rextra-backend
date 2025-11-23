package entity

import (
	"time"

	"github.com/google/uuid"
)

type Memberships struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID           uuid.UUID    `json:"user_id" gorm:"not null;unique"`
	MembershipStatus EnumPlanName `json:"membership_status" gorm:"type:varchar(20)"`

	Plan   MembershipPlans `gorm:"foreignKey:PlanID;references:ID"`
	PlanID uuid.UUID       `json:"plan_id"`

	Duration   MembershipDuration `gorm:"foreignKey:DurationID;references:ID"`
	DurationID uuid.UUID          `json:"duration_id"`

	CurrentTokenBalance int       `json:"current_token_balance" gorm:"default:0"`
	CurrentPoinBalance  int       `json:"current_poin_balance" gorm:"default:0"`
	StartedAt           time.Time `json:"started_at"`
	ExpiredAt           time.Time `json:"expired_at"`
	IsActive            bool      `json:"is_active"`
	AutoRenew           bool      `json:"auto_renew" gorm:"default:false"`

	Timestamp
}

func (m *Memberships) TableName() string {
	return "memberships"
}

func NewMembership(userId uuid.UUID, plan *MembershipPlans, duration *MembershipDuration) Memberships {

	return Memberships{
		UserID:           userId,
		MembershipStatus: plan.PlanName,
		PlanID:           plan.ID,
		DurationID:       duration.ID,
		StartedAt:        time.Now().UTC(),
		ExpiredAt:        time.Now().Add(time.Duration(time.Now().Day() * 30)).UTC(),
	}
}

func (m *Memberships) UpdateMembership(planId uuid.UUID, durationId uuid.UUID) {
	m.PlanID = planId
	m.DurationID = durationId
}

func (m *Memberships) UseMembershipToken(token int) {
	m.CurrentTokenBalance -= token
}

func (m *Memberships) UseMembershipPoin(poin int) {
	m.CurrentPoinBalance -= poin
}

func (m *Memberships) GetTokenBalance() int {
	return m.CurrentTokenBalance
}

func (m *Memberships) GetPoinBalance() int {
	return m.CurrentPoinBalance
}

func (m *Memberships) AddMembershipToken(token int) {
	m.CurrentTokenBalance += token
}

func (m *Memberships) AddMembershipPoinBalance(poin int) {
	m.CurrentPoinBalance += poin
}

func (m *Memberships) IsUserMembershipExpired() bool {
	return time.Now().After(m.ExpiredAt)
}
func (m *Memberships) CalculateTotalToken() int {
	baseToken := m.Plan.MonthlyToken * m.Duration.DurationMonth
	bonusToken := float64(baseToken) * (m.Duration.TokenBonusPercentage / 100)
	totalToken := baseToken + int(bonusToken)
	m.CurrentTokenBalance += totalToken
	return totalToken
}

func (m *Memberships) CalculateRextraPoin() int {
	basePoin := m.Plan.MonthlyToken * m.Duration.DurationMonth
	bonusPoin := basePoin * m.Duration.RextraPoinMultiplier
	totalPoin := basePoin + bonusPoin
	m.CurrentPoinBalance += totalPoin
	return totalPoin
}
