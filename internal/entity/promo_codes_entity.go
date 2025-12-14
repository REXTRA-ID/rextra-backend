package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type (
	PromoType    string
	DiscountType string
)

const (
	PromoMembershipType PromoType = "membership"
	PromoTokenType      PromoType = "token"
	PromoBoth                     = "both"
)

const (
	DiscountPercentage  DiscountType = "percentage"
	DiscountFixedAmount DiscountType = "fixed_amount"
)

type PromoCodes struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code                string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Description         string         `gorm:"type:text" json:"description"`
	PromoType           PromoType      `gorm:"type:varchar(20)"`
	DiscountType        DiscountType   `gorm:"type:varchar(20);not null" json:"discount_type"`
	DiscountValue       float64        `gorm:"type:decimal(10,2)" json:"discount_value"`
	ApplicablePlans     datatypes.JSON `gorm:"type:jsonb" json:"applicable_plans"`
	BonusToken          *int           `gorm:"bonus_token"`
	ApplicableDurations datatypes.JSON `gorm:"type:jsonb" json:"applicable_durations"`
	MinTokenPurchase    *int           `json:"min_token_purchase"`
	MaxUsage            *int           `json:"max_usage"`
	MaxUsagePerUser     *int           `json:"max_usage_per_usage"`
	CurrentUsage        int            `gorm:"default:0" json:"current_usage"`
	ValidFrom           time.Time      `gorm:"not null" json:"valid_from"`
	ValidUntil          time.Time      `gorm:"not null" json:"valid_until"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`

	Timestamp
}

func NewPromoCodes(
	code, description, discountType string,
	discountValue float64,
	applicablePlans, applicbleDuration datatypes.JSON,
	bonusToken, minTokenPurchase, maxUsage, maxUsagePerUser *int,
	validFrom, validUntil time.Time,
) PromoCodes {
	return PromoCodes{
		Code:                code,
		Description:         description,
		DiscountType:        DiscountType(discountType),
		DiscountValue:       discountValue,
		ApplicablePlans:     applicablePlans,
		ApplicableDurations: applicbleDuration,
		MaxUsage:            maxUsage,
		MaxUsagePerUser:     maxUsagePerUser,
		BonusToken:          bonusToken,
		MinTokenPurchase:    minTokenPurchase,
		ValidFrom:           validFrom,
		ValidUntil:          validUntil,
	}
}

func (p *PromoCodes) CalculatePriceWithPercentageDiscount(price float64) float64 {
	priceDiscount := (price * p.DiscountValue) / 100
	priceAfterDiscount := price - priceDiscount
	return priceAfterDiscount
}

func (p *PromoCodes) CalculatePriceWithFixedAmountDiscount(price float64) float64 {
	priceAfterDiscount := price - p.DiscountValue
	return priceAfterDiscount
}

func (p *PromoCodes) UsageExceed() bool {
	return p.CurrentUsage > *p.MaxUsage
}

func (p *PromoCodes) IsValid() bool {
	return p.ValidFrom.Before(p.ValidUntil)
}

func (p *PromoCodes) IncreaseCurrentUsage() {
	p.CurrentUsage++
}

func (p *PromoCodes) ApplyPromoApplicablePlans(plan string) (bool, error) {
	var applicablePlans []string
	if err := json.Unmarshal(p.ApplicablePlans, &applicablePlans); err != nil {
		return true, err
	}

	for i := range applicablePlans {
		if applicablePlans[i] == plan {
			return true, nil
		}
	}

	return false, nil
}

func (p *PromoCodes) ApplyPromoApplicaleDuration(duration int) (bool, error) {
	var applicableDurations []int
	if err := json.Unmarshal(p.ApplicableDurations, &applicableDurations); err != nil {
		return false, err
	}

	for i := range applicableDurations {
		if applicableDurations[i] == duration {
			return true, nil
		}
	}
	return false, nil
}

func (p *PromoCodes) TableName() string {
	return "promo_codes"
}
