package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MappingStatus string

const (
	MappingStatusActive   MappingStatus = "aktif"
	MappingStatusInactive MappingStatus = "nonaktif"
)

type DurationAccessMapping struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanDurationID uuid.UUID `json:"plan_duration_id" gorm:"type:uuid;not null;index:idx_dam_unique,unique"`
	EntitlementID  uuid.UUID `json:"entitlement_id" gorm:"type:uuid;not null;index:idx_dam_unique,unique"`

	EntitlementKey  string `json:"entitlement_key" gorm:"type:varchar(200);not null"`
	EntitlementName string `json:"entitlement_name" gorm:"type:varchar(200);not null"`
	Category        string `json:"category" gorm:"type:varchar(100)"`

	UsageLimit int `json:"usage_limit" gorm:"not null;default:0"`

	Status    MappingStatus `json:"status" gorm:"type:varchar(20);not null;default:'aktif'"`
	CreatedAt time.Time     `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	PlanDuration PlanDuration `json:"plan_duration,omitempty" gorm:"foreignKey:PlanDurationID;references:ID"`
	Entitlement  Entitlement  `json:"entitlement,omitempty" gorm:"foreignKey:EntitlementID;references:ID"`
}

func (DurationAccessMapping) TableName() string {
	return "duration_access_mappings"
}

func (d *DurationAccessMapping) BeforeCreate(tx *gorm.DB) error {
	d.CreatedAt = time.Now().UTC()
	return nil
}

func NewDurationAccessMapping(
	planDurationID uuid.UUID,
	entitlementID uuid.UUID,
	entitlementKey, entitlementName, category string,
	usageLimit int,
) DurationAccessMapping {
	return DurationAccessMapping{
		PlanDurationID:  planDurationID,
		EntitlementID:   entitlementID,
		EntitlementKey:  entitlementKey,
		EntitlementName: entitlementName,
		Category:        category,
		UsageLimit:      usageLimit,
		Status:          MappingStatusActive,
	}
}
