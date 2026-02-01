package entity

import (
	"time"

	"github.com/google/uuid"
)

// Test goal types
type TestGoal string

const (
	TestGoalRecommendation TestGoal = "RECOMMENDATION"
	TestGoalFitCheck       TestGoal = "FIT_CHECK"
)

type CareerProfileTestStatus string

type CareerProfileTestSession struct {
	ID                   int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID               uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	SessionToken         string     `json:"session_token" gorm:"size:100;uniqueIndex;not null"`
	PersonaType          string     `json:"persona_type" gorm:"size:50"`
	TestGoal             TestGoal   `json:"test_goal" gorm:"type:varchar(50);not null;default:'RECOMMENDATION'"`
	TargetProfessionID   *int64     `json:"target_profession_id" gorm:"index"`
	UsesIkigai           bool       `json:"uses_ikigai" gorm:"default:false;not null"`
	Status               string     `json:"status" gorm:"size:20;not null"`
	StartedAt            time.Time  `json:"started_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	CompletedAt          *time.Time `json:"completed_at" gorm:"type:timestamptz"`
	RiasecCompletedAt    *time.Time `json:"riasec_completed_at" gorm:"type:timestamptz"`
	IkigaiCompletedAt    *time.Time `json:"ikigai_completed_at" gorm:"type:timestamptz"`
	AlgorithmVersion     *string    `json:"algorithm_version" gorm:"size:50"`
	QuestionSetVersion   *string    `json:"question_set_version" gorm:"size:50"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (CareerProfileTestSession) TableName() string {
	return "careerprofile_test_sessions"
}
