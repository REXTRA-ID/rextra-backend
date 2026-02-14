package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EnumPlanName string

const (
	PLANSTARTER  EnumPlanName = "Starter"
	PLANBASIC    EnumPlanName = "Basic"
	PLANPRO      EnumPlanName = "Pro"
	PLANMAX      EnumPlanName = "Max"
	PLANSTANDARD EnumPlanName = "Standard"
)

type MembershipPlans struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanName         EnumPlanName   `json:"plan_name" gorm:"type:varchar(20)"`
	MonthlyToken     int            `json:"monthly_token"`
	BaseMonthlyPrice float64        `json:"base_monthly_price" gorm:"type:decimal(10, 2)"`
	Description      string         `json:"description"`
	Benefits         datatypes.JSON `json:"benefits" gorm:"type:jsonb"`
	IsActive         bool           `json:"is_active" gorm:"type:boolean;default:true"`

	Timestamp
}

func (m *MembershipPlans) TableName() string {
	return "membership_plans"
}

func NewMembershipPlans(planName string, MonthlyToken int,
	baseMonthlyPrice float64, description string,
	benefits datatypes.JSON) MembershipPlans {

	return MembershipPlans{
		PlanName:         EnumPlanName(planName),
		MonthlyToken:     MonthlyToken,
		BaseMonthlyPrice: baseMonthlyPrice,
		Description:      description,
		Benefits:         benefits,
	}
}

func (m *MembershipPlans) BeforeCreate(tx *gorm.DB) (err error) {
	time := time.Now().UTC()
	m.CreatedAt = time
	m.UpdatedAt = time
	return nil
}
