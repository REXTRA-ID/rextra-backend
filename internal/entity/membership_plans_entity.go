package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	JSON json.RawMessage

	ListBenefits datatypes.JSON

	EnumPlanName string
)

const (
	PLANSTARTER EnumPlanName = "Starter"
	PLANBASIC   EnumPlanName = "Basic"
	PLANPRO     EnumPlanName = "Pro"
	PLANMAX     EnumPlanName = "Max"
)

type MembershipPlans struct {
	ID               uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanName         EnumPlanName    `json:"plan_name" gorm:"type:varchar(20)"`
	MonthlyToken     uint            `json:"monthly_token"`
	BaseMonthlyPrice decimal.Decimal `json:"base_monthly_token" gorm:"type:decimal(10, 2)"`
	Description      string          `json:"description"`
	Benefits         ListBenefits    `json:"benefits" gorm:"type:jsonb"`
	IsActive         bool            `json:"is_active" gorm:"type:boolean;default:true"`

	Timestamp
}

func (m *MembershipPlans) TableName() string {
	return "membership_plans"
}

func (m *MembershipPlans) BeforeCreate(tx *gorm.DB) (err error) {
	time := time.Now().UTC()
	m.CreatedAt = time
	m.UpdatedAt = time
	return nil
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
