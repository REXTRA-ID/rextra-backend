package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SessionToken struct {
	ID     string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID string `json:"user_id" gorm:"not null"`

	Token        string         `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null"`
	IsActive     bool           `json:"is_active" gorm:"default:true;not null"`
	AuthProvider string         `json:"auth_provider" gorm:"not null"`
	DeviceInfo   datatypes.JSON `json:"device_info" gorm:"type:jsonb"`

	Timestamp
}

func (s *SessionToken) TableName() string {
	return "sessions"
}

func (s *SessionToken) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
