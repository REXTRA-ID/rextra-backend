package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserFavoriteProfession struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"          json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"         json:"user_id"`
	ProfessionID int64     `gorm:"not null;index"                   json:"profession_id"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone" json:"updated_at"`

	User       User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"       json:"user,omitempty"`
	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE" json:"profession,omitempty"`
}

func (UserFavoriteProfession) TableName() string {
	return "user_favorite_professions"
}
