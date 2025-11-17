package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PromoCodes struct { // Opsional
	ID                  uuid.UUID
	Code                string
	Description         string
	DiscountType        string
	DiscountValue       string
	ApplicablePlans     datatypes.JSON
	ApplicableDurations datatypes.JSON
	MaxUsage            int
	CurrentUsage        int
	ValidFrom           time.Time
	ValidUntil          time.Time
	IsActive            bool

	Timestamp
}
