package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserRiasec struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Profile      string    `json:"profile" gorm:"not null"` // RIASEC CODE
	NormalizedScores datatypes.JSON `json:"normalized_scores" gorm:"type:jsonb;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null"`
}