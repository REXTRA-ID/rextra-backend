package entity

import (
	"time"

	"gorm.io/datatypes"
)

type IkigaiTotalScore struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	ScoresData    datatypes.JSON `json:"scores_data" gorm:"type:jsonb;not null"`
	CalculatedAt  time.Time      `json:"calculated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (IkigaiTotalScore) TableName() string {
	return "ikigai_total_scores"
}
