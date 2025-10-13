package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type UserIkigai struct {
	ID           	uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       	uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Profile      	string    `json:"profile" gorm:"type:varchar(255);not null"`
	ChartData    	datatypes.JSON `json:"chart_data" gorm:"type:jsonb;not null"`
	Hash         	string    `json:"hash" gorm:"type:varchar(255);not null"`
	Results      	datatypes.JSON `json:"results" gorm:"type:jsonb;not null"`
	RiasecExplanations datatypes.JSON `json:"riasec_explanations" gorm:"type:jsonb;not null"`
	RiasecMapFull  datatypes.JSON `json:"riasec_map_full" gorm:"type:jsonb;not null"`
	CreatedAt    	time.Time `json:"created_at" gorm:"not null"`
}
