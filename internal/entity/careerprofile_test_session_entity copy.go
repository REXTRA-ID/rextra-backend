package entity

import (
	"time"

	"github.com/google/uuid"
)

type CareerProfileTestStatus string

type CareerProfileTestSession struct {
	ID                int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	SessionToken      string     `json:"session_token" gorm:"size:100;uniqueIndex;not null"`
	Status            string     `json:"status" gorm:"size:20;not null"`
	StartedAt         time.Time  `json:"started_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	CompletedAt       *time.Time `json:"completed_at" gorm:"type:timestamptz"`
	RiasecCompletedAt *time.Time `json:"riasec_completed_at" gorm:"type:timestamptz"`
	IkigaiCompletedAt *time.Time `json:"ikigai_completed_at" gorm:"type:timestamptz"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (CareerProfileTestSession) TableName() string {
	return "careerprofile_test_sessions"
}
