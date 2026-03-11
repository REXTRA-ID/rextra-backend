package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type StudentFeedback struct {
	ID                int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	TestCategoryID    int64          `json:"test_category_id" gorm:"not null"`
	EaseOfUseScore    int            `json:"ease_of_use_score" gorm:"not null"`
	RelevanceScore    int            `json:"relevance_score" gorm:"not null"`
	SatisfactionScore int            `json:"satisfaction_score" gorm:"not null"`
	Obstacles         datatypes.JSON `json:"obstacles" gorm:"type:jsonb;not null;default:'[]'"`
	SubmittedAt       time.Time      `json:"submitted_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	User         User               `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TestCategory KenaliDiriCategory `json:"-" gorm:"foreignKey:TestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (StudentFeedback) TableName() string {
	return "student_feedback"
}
