package entity

import (
	"time"
)

type RiasecResult struct {
	ID                    int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID         int64     `json:"test_session_id" gorm:"uniqueIndex;not null"`
	ScoreR                int       `json:"score_r" gorm:"not null"`
	ScoreI                int       `json:"score_i" gorm:"not null"`
	ScoreA                int       `json:"score_a" gorm:"not null"`
	ScoreS                int       `json:"score_s" gorm:"not null"`
	ScoreE                int       `json:"score_e" gorm:"not null"`
	ScoreC                int       `json:"score_c" gorm:"not null"`
	RiasecCodeID          int64     `json:"riasec_code_id" gorm:"not null"`
	RiasecCodeType        string    `json:"riasec_code_type" gorm:"size:20;not null"`
	IsInconsistentProfile bool      `json:"is_inconsistent_profile" gorm:"default:false;not null"`
	CalculatedAt          time.Time `json:"calculated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	RiasecCode               RiasecCode               `json:"-" gorm:"foreignKey:RiasecCodeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (RiasecResult) TableName() string {
	return "riasec_results"
}