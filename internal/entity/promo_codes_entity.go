package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PromoCodes struct { // Opsional
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code                string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Description         string         `gorm:"type:text" json:"description"`
	DiscountType        string         `gorm:"type:varchar(20);not null" json:"discount_type"`
	DiscountValue       float64        `gorm:"type:decimal(10,2)" json:"discount_value"`
	ApplicablePlans     datatypes.JSON `gorm:"type:jsonb" json:"applicable_plans"`
	ApplicableDurations datatypes.JSON `gorm:"type:jsonb" json:"applicable_durations"`
	MaxUsage            *int           `json:"max_usage"`
	CurrentUsage        int            `gorm:"default:0" json:"current_usage"`
	ValidFrom           time.Time      `gorm:"not null" json:"valid_from"`
	ValidUntil          time.Time      `gorm:"not null" json:"valid_until"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`

	Timestamp
}
