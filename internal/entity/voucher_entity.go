package entity

import (
	"github.com/google/uuid"
	"time"
)

type Voucher struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Code      string     `json:"code" gorm:"uniqueIndex;not null"`
	IsUsed    bool       `json:"is_used" gorm:"default:false;not null"`
	UserID    *uuid.UUID `json:"user_id" gorm:"type:uuid"` // Nullable, set when used
	UsedAt    *time.Time `json:"used_at"`                  // Nullable, set when used
	Timestamp
}

func (v *Voucher) TableName() string {
	return "vouchers"
}
