package entity

import (
	"time"

	"gorm.io/datatypes"
)

type ExpertFeedback struct {
	ID                 int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID      int64          `json:"test_session_id" gorm:"not null"`
	TestCategoryID     int64          `json:"test_category_id" gorm:"not null"`
	ExpertName         string         `json:"expert_name" gorm:"size:255;not null"`
	Profession         string         `json:"profession" gorm:"size:255"`
	Degree             string         `json:"degree" gorm:"size:255"`
	Experience         string         `json:"experience" gorm:"size:255"`
	Education          string         `json:"education" gorm:"size:255"`
	University         string         `json:"university" gorm:"size:255"`
	StudyProgram       string         `json:"study_program" gorm:"size:255"`
	CategoryTest       string         `json:"category_test" gorm:"size:100"`
	TopFiveProfessions datatypes.JSON `json:"top_five_professions" gorm:"type:jsonb;not null;default:'[]'"`
	AccuracyScore      int            `json:"accuracy_score" gorm:"not null"`
	LogicScore         int            `json:"logic_score" gorm:"not null"`
	BenefitScore       int            `json:"benefit_score" gorm:"not null"`
	Obstacles          datatypes.JSON `json:"obstacles" gorm:"type:jsonb;not null;default:'[]'"`
	Suggestions        *string        `json:"suggestions" gorm:"type:text"`
	SubmittedAt        time.Time      `json:"submitted_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	TestSession  CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TestCategory KenaliDiriCategory       `json:"-" gorm:"foreignKey:TestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ExpertFeedback) TableName() string {
	return "expert_feedback"
}
