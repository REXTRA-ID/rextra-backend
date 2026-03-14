package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanDuration struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanID         uuid.UUID `json:"plan_id" gorm:"type:uuid;not null;index:idx_plan_duration_unique,unique"`
	DurationMonths int       `json:"duration_months" gorm:"not null;index:idx_plan_duration_unique,unique"` // 1, 3, 6, atau 12

	Price       int64   `json:"price" gorm:"not null;default:0"`
	DiscountPct float64 `json:"discount_pct" gorm:"type:decimal(5,2);default:0"`
	FinalPrice  int64   `json:"final_price" gorm:"not null;default:0"`

	DurationPrice int64 `json:"duration_price" gorm:"not null;default:0"`

	TokenAmount int `json:"token_amount" gorm:"not null;default:0"`
	BonusToken  int `json:"bonus_token" gorm:"not null;default:0"`

	PointsActive bool `json:"points_active" gorm:"not null;default:false"`
	PointsValue  int  `json:"points_value" gorm:"not null;default:0"`
	BonusPoints  int  `json:"bonus_points" gorm:"not null;default:0"`

	IsActive bool `json:"is_active" gorm:"not null;default:true"`

	Timestamp

	Plan *MembershipPlans `json:"plan,omitempty" gorm:"foreignKey:PlanID;references:ID"`

	DurationAccessMappings []DurationAccessMapping `json:"duration_access_mappings,omitempty" gorm:"foreignKey:PlanDurationID;references:ID"`

	Memberships []Memberships `json:"memberships,omitempty" gorm:"foreignKey:DurationID;references:ID"`

	SubscriptionCycles []SubscriptionCycle `json:"subscription_cycles,omitempty" gorm:"foreignKey:PlanDurationID;references:ID"`
}

func (PlanDuration) TableName() string {
	return "plan_durations"
}

func (p *PlanDuration) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (p *PlanDuration) TotalDays() int {
	return p.DurationMonths * 30
}

func (p *PlanDuration) DailyRate() float64 {
	if p.TotalDays() == 0 {
		return 0
	}
	return float64(p.DurationPrice) / float64(p.TotalDays())
}

func (p *PlanDuration) CalculateCredit(remainingDays int) int64 {
	return int64(p.DailyRate() * float64(remainingDays))
}

func NewPlanDuration(
	planID uuid.UUID,
	durationMonths int,
	price, finalPrice, durationPrice int64,
	discountPct float64,
	tokenAmount, bonusToken int,
) PlanDuration {
	return PlanDuration{
		PlanID:         planID,
		DurationMonths: durationMonths,
		Price:          price,
		DiscountPct:    discountPct,
		FinalPrice:     finalPrice,
		DurationPrice:  durationPrice,
		TokenAmount:    tokenAmount,
		BonusToken:     bonusToken,
		IsActive:       true,
	}
}
