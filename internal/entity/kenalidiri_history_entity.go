package entity

import (
	"time"

	"github.com/google/uuid"
)

type KenaliDiriHistoryStatus string

type KenaliDiriHistory struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	TestCategoryID  int64      `json:"test_category_id" gorm:"not null;index"`
	DetailSessionID int64      `json:"detail_session_id" gorm:"not null"`
	Status          string     `json:"status" gorm:"size:20;not null"`
	StartedAt       time.Time  `json:"started_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	CompletedAt     *time.Time `json:"completed_at" gorm:"type:timestamptz"`

	User         User               `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TestCategory KenaliDiriCategory `json:"-" gorm:"foreignKey:TestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (KenaliDiriHistory) TableName() string {
	return "kenalidiri_history"
}
