package entity

import (
	"time"

	"gorm.io/datatypes"
)

type CareerRecommendation struct {
	ID                  int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID       int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	RecommendationsData datatypes.JSON `json:"recommendations_data" gorm:"type:jsonb;not null"`
	TopProfession1ID    *int64         `json:"top_profession1_id" gorm:"type:bigint"`
	TopProfession2ID    *int64         `json:"top_profession2_id" gorm:"type:bigint"`
	GeneratedAt         time.Time      `json:"generated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	AIModelUsed         string         `json:"ai_model_used" gorm:"size:50;default:'gemini-1.5-flash';not null"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CareerRecommendation) TableName() string {
	return "career_recommendations"
}
