package entity

import (
	"time"

	"github.com/google/uuid"
)

type RespondentType string

const (
	RespondentTypeStudent RespondentType = "STUDENT"
	RespondentTypeExpert  RespondentType = "EXPERT"
)

type KenaliDiriFeedback struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestCategory     string         `json:"test_category" gorm:"size:50;not null;index"`
	TestSessionID    int64          `json:"test_session_id" gorm:"not null;index"`
	RespondentType   RespondentType `json:"respondent_type" gorm:"type:varchar(20);not null;index"`
	RespondentUserID uuid.UUID      `json:"respondent_user_id" gorm:"type:uuid;not null;index"`
	SubmittedAt      time.Time      `json:"submitted_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	DeletedAt        *time.Time     `json:"deleted_at" gorm:"type:timestamptz;index"`
	DeletedBy        *uuid.UUID     `json:"deleted_by" gorm:"type:uuid"`

	RespondentUser User  `json:"-" gorm:"foreignKey:RespondentUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeletedByUser  *User `json:"-" gorm:"foreignKey:DeletedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (KenaliDiriFeedback) TableName() string {
	return "kenalidiri_feedback"
}
