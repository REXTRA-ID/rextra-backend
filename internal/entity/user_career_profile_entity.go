package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserCareerProfile struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          uuid.UUID  `json:"user_id" gorm:"uniqueIndex;not null"`
	CurrentPlan     string     `json:"current_plan" gorm:"size:50;default:'Standard';not null"`
	ProfileData     string     `json:"profile_data" gorm:"type:text"`
	ActiveSessionID *int64     `json:"active_session_id" gorm:"type:bigint"`
	SetSource       string     `json:"set_source" gorm:"size:100"`
	SetAt           *time.Time `json:"set_at" gorm:"type:timestamptz"`
	Pinned          bool       `json:"pinned" gorm:"default:false"`
	LastUpdateAt    time.Time  `json:"last_update_at" gorm:"type:timestamptz;default:now();autoUpdateTime"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserCareerProfile) TableName() string {
	return "user_career_profiles"
}
