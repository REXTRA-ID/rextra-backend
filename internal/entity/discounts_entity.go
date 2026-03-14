package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DiscountType string
type DiscountAppliesTo string
type DiscountStatus string

const (
	DiscountTypePercentage DiscountType = "PERCENTAGE"
	DiscountTypeFixed      DiscountType = "FIXED"

	DiscountAppliesToMembership DiscountAppliesTo = "MEMBERSHIP"
	DiscountAppliesToTokenTopup DiscountAppliesTo = "TOKEN_TOPUP"
	DiscountAppliesToGlobal     DiscountAppliesTo = "GLOBAL"

	DiscountStatusActive   DiscountStatus = "ACTIVE"
	DiscountStatusInactive DiscountStatus = "INACTIVE"
)

type Discounts struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name string    `json:"name" gorm:"type:varchar(100);not null"`

	DiscountType DiscountType      `json:"discount_type" gorm:"type:varchar(20);not null"`
	Value        float64           `json:"value" gorm:"type:decimal(12,2);not null"`
	AppliesTo    DiscountAppliesTo `json:"applies_to" gorm:"type:varchar(20);not null;default:'MEMBERSHIP'"`

	MembershipPlanTargets datatypes.JSON `json:"membership_plan_targets,omitempty" gorm:"type:jsonb"`

	TopupTargets datatypes.JSON `json:"topup_targets,omitempty" gorm:"type:jsonb"`

	MaxDiscountAmount *int64 `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount *int64 `json:"min_purchase_amount,omitempty"`

	MaxTotalRedemptions   *int `json:"max_total_redemptions,omitempty"`
	MaxRedemptionsPerUser *int `json:"max_redemptions_per_user,omitempty"`
	CurrentRedemptions    int  `json:"current_redemptions" gorm:"not null;default:0"`

	Priority  int  `json:"priority" gorm:"not null;default:0"`
	Stackable bool `json:"stackable" gorm:"not null;default:false"`

	StartsAt *time.Time `json:"starts_at,omitempty"`
	EndsAt   *time.Time `json:"ends_at,omitempty"`

	Description string         `json:"description" gorm:"type:text"`
	Status      DiscountStatus `json:"status" gorm:"type:varchar(20);not null;default:'ACTIVE'"`

	CreatedBy string `json:"created_by" gorm:"type:varchar(100);default:'System'"`
	UpdatedBy string `json:"updated_by" gorm:"type:varchar(100);default:'System'"`

	Timestamp

	Redemptions []DiscountRedemption `json:"redemptions,omitempty" gorm:"foreignKey:DiscountID;references:ID"`
}

func (Discounts) TableName() string {
	return "discounts"
}

func (d *Discounts) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	return nil
}

func (d *Discounts) IsValid() bool {
	if d.Status != DiscountStatusActive {
		return false
	}
	now := time.Now().UTC()
	if d.StartsAt != nil && now.Before(*d.StartsAt) {
		return false
	}
	if d.EndsAt != nil && now.After(*d.EndsAt) {
		return false
	}
	if d.MaxTotalRedemptions != nil && d.CurrentRedemptions >= *d.MaxTotalRedemptions {
		return false
	}
	return true
}

func (d *Discounts) CalculateDiscount(price int64) int64 {
	var discount int64
	if d.DiscountType == DiscountTypePercentage {
		discount = int64(float64(price) * d.Value / 100)
		if d.MaxDiscountAmount != nil && discount > *d.MaxDiscountAmount {
			discount = *d.MaxDiscountAmount
		}
	} else {
		discount = int64(d.Value)
	}
	if discount > price {
		return price
	}
	return discount
}

func (d *Discounts) IncrementRedemption() {
	d.CurrentRedemptions++
}

func (d *Discounts) DecrementRedemption() {
	if d.CurrentRedemptions > 0 {
		d.CurrentRedemptions--
	}
}

func NewDiscount(
	code, name string,
	discountType DiscountType,
	value float64,
	appliesTo DiscountAppliesTo,
	priority int,
	stackable bool,
	createdBy string,
) Discounts {
	return Discounts{
		Code:         code,
		Name:         name,
		DiscountType: discountType,
		Value:        value,
		AppliesTo:    appliesTo,
		Priority:     priority,
		Stackable:    stackable,
		Status:       DiscountStatusActive,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
	}
}
