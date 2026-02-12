package entity

import (
	"time"

	"github.com/google/uuid"
)

type TokenWallet struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;column:user_id"`
	Balance   int64     `json:"balance" gorm:"not null;default:0;column:balance"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;column:updated_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

func (TokenWallet) TableName() string {
	return "token_wallet"
}
