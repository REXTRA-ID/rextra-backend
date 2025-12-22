package entity

import (
	"time"

	"gorm.io/datatypes"
)

type IkigaiDimensionScore struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	ScoresData    datatypes.JSON `json:"scores_data" gorm:"type:jsonb;not null"`
	CalculatedAt  time.Time      `json:"calculated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	AIModelUsed   string         `json:"ai_model_used" gorm:"size:50;default:'gemini-1.5-flash';not null"`
	TotalAPICalls int            `json:"total_api_calls" gorm:"default:4"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
