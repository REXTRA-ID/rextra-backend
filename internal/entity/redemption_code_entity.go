package entity

import (
	"time"

	"github.com/google/uuid"
)

type RedemptionCode struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code       string    `json:"code" gorm:"not null"`
	TokenValue int       `json:"token_value" gorm:"not null"`
	IsUsed     bool      `json:"is_used" gorm:"default:false;not null"`
	RedeemedBy uuid.UUID `json:"redeemed_by" gorm:""`
	RedeemedAt time.Time `json:"redeemed_at" gorm:""`
}

func (r *RedemptionCode) TableName() string {
	return "redemption_codes"
}
