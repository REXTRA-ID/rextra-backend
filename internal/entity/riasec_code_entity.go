package entity

import (
	"time"

	"gorm.io/datatypes"
)

type RiasecCode struct {
	ID                int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	RiasecCode        string         `json:"riasec_code" gorm:"size:3;uniqueIndex;not null"`
	RiasecTitle       string         `json:"riasec_title" gorm:"size:255;not null"`
	RiasecDescription *string        `json:"riasec_description" gorm:"type:text"`
	Strengths         datatypes.JSON `json:"strengths" gorm:"type:jsonb;not null;default:'[]'"`
	Challenges        datatypes.JSON `json:"challenges" gorm:"type:jsonb;not null;default:'[]'"`
	Strategies        datatypes.JSON `json:"strategies" gorm:"type:jsonb;not null;default:'[]'"`
	WorkEnvironments  datatypes.JSON `json:"work_environments" gorm:"type:jsonb;not null;default:'[]'"`
	InteractionStyles datatypes.JSON `json:"interaction_styles" gorm:"type:jsonb;not null;default:'[]'"`
	CreatedAt         time.Time      `json:"created_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"type:timestamptz;default:now();autoUpdateTime"`
}

func (RiasecCode) TableName() string {
	return "riasec_codes"
}
