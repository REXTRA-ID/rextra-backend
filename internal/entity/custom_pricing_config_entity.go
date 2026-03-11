package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomPricingMetadata for JSONB field
type CustomPricingMetadata map[string]any

// Scan implements sql.Scanner
func (m *CustomPricingMetadata) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

// Value implements driver.Valuer
func (m CustomPricingMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type CustomPricingConfig struct {
	ID                       uuid.UUID             `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IsEnabled                bool                  `json:"is_enabled" gorm:"not null;default:false;column:is_enabled"`
	MinToken                 int64                 `json:"min_token" gorm:"not null;column:min_token"`
	MaxToken                 int64                 `json:"max_token" gorm:"not null;column:max_token"`
	RecommendedPricePerToken int64                 `json:"recommended_price_per_token" gorm:"not null;column:recommended_price_per_token"`
	GuardrailAckBy           *uuid.UUID            `json:"guardrail_ack_by,omitempty" gorm:"type:uuid;column:guardrail_ack_by"`
	GuardrailAckAt           *time.Time            `json:"guardrail_ack_at,omitempty" gorm:"column:guardrail_ack_at"`
	IsCurrent                bool                  `json:"is_current" gorm:"not null;default:false;column:is_current"`
	EffectiveFrom            time.Time             `json:"effective_from" gorm:"not null;default:CURRENT_TIMESTAMP;column:effective_from"`
	EffectiveTo              *time.Time            `json:"effective_to,omitempty" gorm:"column:effective_to"`
	UpdatedBy                *uuid.UUID            `json:"updated_by,omitempty" gorm:"type:uuid;column:updated_by"`
	Metadata                 CustomPricingMetadata `json:"metadata" gorm:"type:jsonb;not null;default:'{}';column:metadata"`
	Timestamp

	// Relations
	Tiers          []CustomPricingTier `json:"tiers,omitempty" gorm:"foreignKey:ConfigID;references:ID"`
	UpdatedByUser  *User               `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBy;references:ID"`
	GuardrailAcker *User               `json:"guardrail_acker,omitempty" gorm:"foreignKey:GuardrailAckBy;references:ID"`
}

// TableName override
func (CustomPricingConfig) TableName() string {
	return "custom_pricing_config"
}

func (c *CustomPricingConfig) BeforeSave(tx *gorm.DB) error {
	if c.MinToken < 1 {
		return errors.New("min_token must be at least 1")
	}

	if c.MaxToken < 1 {
		return errors.New("max_token must be at least 1")
	}

	if c.MinToken >= c.MaxToken {
		return errors.New("min_token must be less than max_token")
	}

	if c.RecommendedPricePerToken < 1 {
		return errors.New("recommended_price_per_token must be at least 1")
	}

	return nil
}
