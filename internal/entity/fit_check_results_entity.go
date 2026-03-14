package entity

import (
	"time"
)

type MatchCategory string

const (
	MatchCategoryHigh   MatchCategory = "HIGH"
	MatchCategoryMedium MatchCategory = "MEDIUM"
	MatchCategoryLow    MatchCategory = "LOW"
)

type FitCheckResult struct {
	ID                      int64         json:"id" gorm:"primaryKey;autoIncrement"
	TestSessionID           int64         json:"test_session_id" gorm:"uniqueIndex;not null"
	ProfessionID            int64         json:"profession_id" gorm:"not null"
	UserRIASECCodeID        int64         json:"user_riasec_code_id" gorm:"not null"
	ProfessionRIASECCodeID  int64         json:"profession_riasec_code_id" gorm:"not null"
	MatchCategory           MatchCategory json:"match_category" gorm:"type:match_category_enum;not null"
	RuleType                string        json:"rule_type" gorm:"type:varchar(50);not null"
	DominantLetterSame      bool          json:"dominant_letter_same" gorm:"not null"
	IsAdjacentHexagon       bool          json:"is_adjacent_hexagon" gorm:"not null"
	MatchScore              float64       json:"match_score" gorm:"type:numeric(4,2)"
	CreatedAt               time.Time     json:"created_at" gorm:"type:timestamptz;default:now();not null"

	CareerProfileTestSession CareerProfileTestSession json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"
}

func (FitCheckResult) TableName() string {
	return "fit_check_results"
}
