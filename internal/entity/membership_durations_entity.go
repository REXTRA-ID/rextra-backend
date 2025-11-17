package entity

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
	"gorm.io/gorm"
)

type MembershipDuration struct {
	ID                   uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DurationMonth        int             `json:"duration_months"`
	TokenBonusPercentage decimal.Decimal `json:"token_bonus_percentage"`
	RextraPoinMultiplier int             `json:"rextra_point_multiplier"`
	IsActive             bool            `json:"is_active" gorm:"type:boolean;default:true"`

	Timestamp
}

func (m *MembershipDuration) BeforeCreate(tx *gorm.DB) (err error) {
	time := time.Now().UTC()
	m.CreatedAt = time
	m.UpdatedAt = time
	return nil
}
