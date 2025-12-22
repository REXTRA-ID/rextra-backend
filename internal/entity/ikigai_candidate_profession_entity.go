package entity

import (
	"time"

	"gorm.io/datatypes"
)

type IkigaiCandidateProfession struct {
	ID                 int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TestSessionID      int64          `json:"test_session_id" gorm:"uniqueIndex;not null"`
	CandidatesData     datatypes.JSON `json:"candidates_data" gorm:"type:jsonb;not null"`
	TotalCandidates    int            `json:"total_candidates" gorm:"not null"`
	GenerationStrategy string         `json:"generation_strategy" gorm:"size:50;default:'4_tier_expansion';not null"`
	MaxCandidatesLimit int            `json:"max_candidates_limit" gorm:"default:15;not null"`
	GeneratedAt        time.Time      `json:"generated_at" gorm:"type:timestamptz;default:now();autoCreateTime"`

	CareerProfileTestSession CareerProfileTestSession `json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (IkigaiCandidateProfession) TableName() string {
	return "ikigai_candidate_professions"
}
