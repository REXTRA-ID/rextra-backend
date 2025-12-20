package entity

import (
	"time"

	"gorm.io/datatypes"
)

type RiasecQuestionSet struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	QuestionIDs   datatypes.JSON `json:"question_ids" gorm:"type:jsonb;not null"`
	GeneratedAt   time.Time      `json:"generated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (RiasecQuestionSet) TableName() string {
	return "riasec_question_sets"
}