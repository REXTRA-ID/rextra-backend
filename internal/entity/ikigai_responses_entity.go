package entity

import (
	"time"

	"gorm.io/datatypes"
)

type IkigaiResponse struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	ResponsesData datatypes.JSON `json:"responses_data" gorm:"type:jsonb;not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (IkigaiResponse) TableName() string {
	return "ikigai_responses"
}