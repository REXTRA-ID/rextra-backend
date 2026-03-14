package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Memberships struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`

	PlanID   *uuid.UUID `json:"plan_id,omitempty" gorm:"type:uuid;index"`
	PlanName PlanName   `json:"plan_name" gorm:"type:varchar(20);not null"`

	DurationID     *uuid.UUID `json:"duration_id,omitempty" gorm:"type:uuid;index"`
	DurationMonths *int       `json:"duration_months,omitempty"`

	StartedAt *time.Time `json:"started_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`

	IsActive  bool `json:"is_active" gorm:"not null;default:true"`
	AutoRenew bool `json:"auto_renew" gorm:"not null;default:false"`

	CurrentTokenBalance int `json:"current_token_balance" gorm:"not null;default:0"`
	CurrentPoinBalance  int `json:"current_poin_balance" gorm:"not null;default:0"`

	PaidCycleCount   int `json:"paid_cycle_count" gorm:"not null;default:0"`
	EntitlementCount int `json:"entitlement_count" gorm:"not null;default:0"`

	Timestamp

	Plan *MembershipPlans `json:"plan,omitempty" gorm:"foreignKey:PlanID;references:ID"`

	Duration *PlanDuration `json:"duration,omitempty" gorm:"foreignKey:DurationID;references:ID"`

	SubscriptionCycles []SubscriptionCycle `json:"subscription_cycles,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	PaymentTransactions []PaymentTransactions `json:"payment_transactions,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	PoinTransactions []PoinTransactions `json:"poin_transactions,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	PointsLedger []PointsLedger `json:"points_ledger,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	UsageLogs []UsageLog `json:"usage_logs,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	UserEntitlementQuotas []UserEntitlementQuota `json:"user_entitlement_quotas,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (Memberships) TableName() string {
	return "memberships"
}

func (m *Memberships) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func NewMembership(userID uuid.UUID, planName PlanName) Memberships {
	return Memberships{
		UserID:   userID,
		PlanName: planName,
		IsActive: true,
	}
}

func (m *Memberships) Activate(plan *MembershipPlans, duration *PlanDuration) {
	now := time.Now().UTC()
	expired := now.AddDate(0, duration.DurationMonths, 0)

	m.PlanID = &plan.ID
	m.PlanName = plan.PlanName
	m.DurationID = &duration.ID
	m.DurationMonths = &duration.DurationMonths
	m.StartedAt = &now
	m.ExpiredAt = &expired
	m.IsActive = true
	m.PaidCycleCount++
	m.UpdatedAt = now
}

func (m *Memberships) Deactivate(fallbackPlanID uuid.UUID, fallbackPlanName PlanName) {
	now := time.Now().UTC()
	m.PlanID = &fallbackPlanID
	m.PlanName = fallbackPlanName
	m.DurationID = nil
	m.DurationMonths = nil
	m.StartedAt = nil
	m.ExpiredAt = nil
	m.IsActive = false
	m.UpdatedAt = now
}

func (m *Memberships) RemainingDays() int {
	if m.ExpiredAt == nil {
		return 0
	}
	remaining := time.Until(*m.ExpiredAt).Hours() / 24
	if remaining < 0 {
		return 0
	}
	return int(remaining)
}

func (m *Memberships) IsEligibleForChangeOfPlan() bool {
	if !m.IsActive || m.ExpiredAt == nil {
		return false
	}
	return time.Now().Before(*m.ExpiredAt)
}

func (m *Memberships) AddPoinBalance(amount int) {
	m.CurrentPoinBalance += amount
	m.UpdatedAt = time.Now().UTC()
}

func (m *Memberships) DeductPoinBalance(amount int) {
	m.CurrentPoinBalance -= amount
	if m.CurrentPoinBalance < 0 {
		m.CurrentPoinBalance = 0
	}
	m.UpdatedAt = time.Now().UTC()
}

func (m *Memberships) UseMembershipToken(amount int) {
	m.CurrentTokenBalance -= amount
	if m.CurrentTokenBalance < 0 {
		m.CurrentTokenBalance = 0
	}
	m.UpdatedAt = time.Now().UTC()
}

func (m *Memberships) RefillToken(amount int) {
	m.CurrentTokenBalance = amount
	m.UpdatedAt = time.Now().UTC()
}

func (m *Memberships) ExpiringMembership() {
	m.IsActive = false
	m.UpdatedAt = time.Now().UTC()
}
