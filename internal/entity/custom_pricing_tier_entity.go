package entity

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomPricingTier represents discount tier
type CustomPricingTier struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ConfigID    uuid.UUID `json:"config_id" gorm:"type:uuid;not null;column:config_id;index:idx_tier_config_from"`
	FromToken   int64     `json:"from_token" gorm:"not null;column:from_token"`
	ToToken     int64     `json:"to_token" gorm:"not null;column:to_token"`
	DiscountPct float64   `json:"discount_pct" gorm:"type:numeric(5,2);not null;column:discount_pct"`
	Timestamp

	// Relations
	Config CustomPricingConfig `json:"-" gorm:"foreignKey:ConfigID;references:ID"`

	// Computed field (not stored in DB)
	EffectivePricePerToken float64 `json:"effective_price_per_token,omitempty" gorm:"-"`
}

// TableName override
func (CustomPricingTier) TableName() string {
	return "custom_pricing_tier"
}

// BeforeSave validation
func (t *CustomPricingTier) BeforeSave(tx *gorm.DB) error {
	if t.FromToken < 1 {
		return errors.New("from_token must be at least 1")
	}

	if t.ToToken < 1 {
		return errors.New("to_token must be at least 1")
	}

	if t.FromToken > t.ToToken {
		return errors.New("from_token must be <= to_token")
	}

	if t.DiscountPct < 0 || t.DiscountPct > 100 {
		return errors.New("discount_pct must be between 0 and 100")
	}

	return nil
}
