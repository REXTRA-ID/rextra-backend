package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PlanName string
type PlanCategory string
type PlanStatus string
type PricingMode string
type DurationMode string

const (
	PlanNameStandard PlanName = "Standard"
	PlanNameStarter  PlanName = "Starter"
	PlanNameBasic    PlanName = "Basic"
	PlanNamePro      PlanName = "Pro"
	PlanNameMax      PlanName = "Max"

	PlanCategoryUnpaid PlanCategory = "unpaid"
	PlanCategoryPaid   PlanCategory = "paid"

	PlanStatusActive   PlanStatus = "aktif"
	PlanStatusInactive PlanStatus = "nonaktif"

	PricingModeManual    PricingMode = "manual"
	PricingModeAutomatic PricingMode = "otomatis"

	DurationModeWithDuration    DurationMode = "dengan_durasi"
	DurationModeWithoutDuration DurationMode = "tanpa_durasi"
)

type MembershipPlans struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanName     PlanName     `json:"plan_name" gorm:"type:varchar(20);uniqueIndex;not null"`
	Category     PlanCategory `json:"category" gorm:"type:varchar(10);not null"`
	TierLabel    string       `json:"tier_label" gorm:"type:varchar(50);not null"`
	EmblemKey    string       `json:"emblem_key" gorm:"type:varchar(50)"`
	Description  string       `json:"description" gorm:"type:text"`
	Status       PlanStatus   `json:"status" gorm:"type:varchar(20);not null;default:'aktif'"`
	PricingMode  PricingMode  `json:"pricing_mode" gorm:"type:varchar(20);not null;default:'manual'"`
	DurationMode DurationMode `json:"duration_mode" gorm:"type:varchar(20);not null;default:'dengan_durasi'"`

	BasePrice1M int64 `json:"base_price_1m" gorm:"not null;default:0"`
	BaseToken1M int   `json:"base_token_1m" gorm:"not null;default:0"`

	Discount3M  float64 `json:"discount_3m" gorm:"type:decimal(5,2);default:0"`
	Discount6M  float64 `json:"discount_6m" gorm:"type:decimal(5,2);default:0"`
	Discount12M float64 `json:"discount_12m" gorm:"type:decimal(5,2);default:0"`

	BonusToken3M  int `json:"bonus_token_3m" gorm:"default:0"`
	BonusToken6M  int `json:"bonus_token_6m" gorm:"default:0"`
	BonusToken12M int `json:"bonus_token_12m" gorm:"default:0"`

	Benefits    datatypes.JSON `json:"benefits" gorm:"type:jsonb;default:'[]'"`
	ActiveUsers int            `json:"active_users" gorm:"default:0"`

	StarterDurationMonths int `json:"starter_duration_months" gorm:"not null;default:1"`

	Timestamp

	PlanDurations       []PlanDuration        `json:"plan_durations,omitempty" gorm:"foreignKey:PlanID;references:ID"`
	Memberships         []Memberships         `json:"memberships,omitempty" gorm:"foreignKey:PlanID;references:ID"`
	PaymentTransactions []PaymentTransactions `json:"payment_transactions,omitempty" gorm:"foreignKey:PlanID;references:ID"`
}

func (MembershipPlans) TableName() string { return "membership_plans" }

func (m *MembershipPlans) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func (m *MembershipPlans) GetBenefits() ([]string, error) {
	var benefits []string
	if err := json.Unmarshal(m.Benefits, &benefits); err != nil {
		return nil, err
	}
	return benefits, nil
}

func NewMembershipPlan(
	planName PlanName,
	category PlanCategory,
	tierLabel, emblemKey, description string,
	pricingMode PricingMode,
	durationMode DurationMode,
	basePrice1M int64,
	baseToken1M int,
) MembershipPlans {
	return MembershipPlans{
		PlanName:              planName,
		Category:              category,
		TierLabel:             tierLabel,
		EmblemKey:             emblemKey,
		Description:           description,
		Status:                PlanStatusActive,
		PricingMode:           pricingMode,
		DurationMode:          durationMode,
		BasePrice1M:           basePrice1M,
		BaseToken1M:           baseToken1M,
		StarterDurationMonths: 1,
		Benefits:              datatypes.JSON([]byte("[]")),
	}
}
