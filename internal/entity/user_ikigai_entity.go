package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type UserIkigai struct {
	ID           	uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       	uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Analysis     	datatypes.JSON `json:"analysis" gorm:"type:jsonb;not null"`
	TopProfessions 	pq.StringArray `json:"top_professions" gorm:"type:text[];not null"`
	CreatedAt    	time.Time `json:"created_at" gorm:"not null"`
}