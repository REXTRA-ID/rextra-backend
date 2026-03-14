package entity

import (
	"time"
)

type IkigaiTotalScores struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID int64     `json:"test_session_id" gorm:"uniqueIndex;not null"`
	PassionScore  int       `json:"passion_score" gorm:"not null"`
	MissionScore  int       `json:"mission_score" gorm:"not null"`
	VocationScore int       `json:"vocation_score" gorm:"not null"`
	ProfessionScore int     `json:"profession_score" gorm:"not null"`
	TotalIkigai   int       `json:"total_ikigai" gorm:"not null"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamptz;default:now();autoUpdateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type IkigaiTotalScore = IkigaiTotalScores

func (IkigaiTotalScores) TableName() string {
	return "ikigai_total_scores"
}
