package entity

import (
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null"`
	TokenBalance int       `json:"token_balance" gorm:"not null"`
	ExpiredAt    time.Time `json:"expired_at" gorm:"not null"`
	Timestamp
}

func (m *Membership) TableName() string {
	return "memberships"
}
