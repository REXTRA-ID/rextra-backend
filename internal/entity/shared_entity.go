package entity

import (
	"time"

	"gorm.io/gorm"
)

type TokenDirection string
type TokenSourceType string

// Timestamp di-embed oleh hampir semua entity di sistem ini.
type Timestamp struct {
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

const (
	PLANSTARTER  = "Starter"
	PLANSTANDARD = "Standard"
	PLANBASIC    = "Basic"
	PLANPRO      = "Pro"
	PLANMAX      = "Max"

	KENALIDIRI    = "kenali_diri"
	CVGENERATOR   = "cv_generator"
	AIINTERVIEWER = "ai_interviewer"

	TokenDirectionIN  TokenDirection = "IN"
	TokenDirectionOUT TokenDirection = "OUT"

	TokenSourceMembership TokenSourceType = "MEMBERSHIP"
	TokenSourceUsage      TokenSourceType = "USAGE"

	TOKENUSAGE         = "usage"
	TOKENMONTHLYREFILL = "monthly_refill"

	SetSourceAutoFirstTime = "AUTO_FIRST_TIME"
)
