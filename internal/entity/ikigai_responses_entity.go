package entity

import (
	"time"

	"gorm.io/datatypes"
)

type IkigaiResponse struct {
	ID                   int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID        int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	Dimension1Love       datatypes.JSON `json:"dimension_1_love" gorm:"type:jsonb;not null"`
	Dimension2GoodAt     datatypes.JSON `json:"dimension_2_good_at" gorm:"type:jsonb;not null"`
	Dimension3WorldNeeds datatypes.JSON `json:"dimension_3_world_needs" gorm:"type:jsonb;not null"`
	Dimension4PaidFor    datatypes.JSON `json:"dimension_4_paid_for" gorm:"type:jsonb;not null"`
	Completed            bool           `json:"completed" gorm:"not null;default:false"`
	CreatedAt            time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	CompletedAt          *time.Time     `json:"completed_at" gorm:"type:timestamptz;check:valid_completion,((completed = FALSE AND completed_at IS NULL) OR (completed = TRUE AND completed_at IS NOT NULL))"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (IkigaiResponse) TableName() string {
	return "ikigai_responses"
}
