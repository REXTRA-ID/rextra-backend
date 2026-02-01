package entity

import (
	"time"

	"github.com/google/uuid"
)

type SetSource string

const (
	SetSourceAutoFirstTime SetSource = "AUTO_FIRST_TIME"
	SetSourceUserSelected  SetSource = "USER_SELECTED"
)

type UserCareerProfile struct {
	UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	ActiveSessionID int64     `json:"active_session_id" gorm:"not null;index"`
	Pinned          bool      `json:"pinned" gorm:"default:true;not null"`
	SetSource       SetSource `json:"set_source" gorm:"type:varchar(50);not null;default:'AUTO_FIRST_TIME'"`
	SetAt           time.Time `json:"set_at" gorm:"type:timestamptz;default:now();not null"`

	User          User                     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActiveSession CareerProfileTestSession `json:"-" gorm:"foreignKey:ActiveSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (UserCareerProfile) TableName() string {
	return "user_career_profiles"
}
