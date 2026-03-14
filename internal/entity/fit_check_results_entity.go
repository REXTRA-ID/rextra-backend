package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FitCheckResults struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64          `json:"test_session_id" gorm:"index;not null"`
	TargetUserID  uuid.UUID      `json:"target_user_id" gorm:"type:uuid;not null"`
	ProfessionID  int64          `json:"profession_id" gorm:"not null"`
	ProfessionName string        `json:"profession_name" gorm:"size:255;not null"`
	MatchScore    int            `json:"match_score" gorm:"not null"`
	AnalysisData  datatypes.JSON `json:"analysis_data" gorm:"type:jsonb;not null"`
	GeneratedAt   time.Time      `json:"generated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (FitCheckResults) TableName() string {
	return "fit_check_results"
}
