package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MappingStatus string
type RestrictionType string

const (
	MappingStatusActive   MappingStatus = "aktif"
	MappingStatusInactive MappingStatus = "nonaktif"

	RestrictionUnlimited        RestrictionType = "unlimited"
	RestrictionTokenGated       RestrictionType = "token_gated"
	RestrictionFrequencyLimited RestrictionType = "frequency_limited"
	RestrictionLocked           RestrictionType = "locked"
)

type DurationAccessMapping struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PlanDurationID uuid.UUID `json:"plan_duration_id" gorm:"type:uuid;not null;index:idx_dam_unique,unique"`
	EntitlementID  uuid.UUID `json:"entitlement_id" gorm:"type:uuid;not null;index:idx_dam_unique,unique"`

	RestrictionType RestrictionType `json:"restriction_type" gorm:"type:varchar(30);not null;default:'unlimited'"`
	TokenCost       int             `json:"token_cost" gorm:"not null;default:0"`
	UsageLimit      int             `json:"usage_limit" gorm:"not null;default:0"`
	ResetPeriod     *string         `json:"reset_period,omitempty" gorm:"type:varchar(20)"`

	// snapshot untuk admin UI
	EntitlementKey  string `json:"entitlement_key" gorm:"type:varchar(200);not null"`
	EntitlementName string `json:"entitlement_name" gorm:"type:varchar(200);not null"`
	Category        string `json:"category" gorm:"type:varchar(100)"`

	Status    MappingStatus `json:"status" gorm:"type:varchar(20);not null;default:'aktif'"`
	CreatedAt time.Time     `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// relasi ke membership_durations
	PlanDuration MembershipDuration `json:"plan_duration,omitempty" gorm:"foreignKey:PlanDurationID;references:ID"`

	// relasi ke entitlement
	Entitlement Entitlement `json:"entitlement,omitempty" gorm:"foreignKey:EntitlementID;references:ID"`
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
	entitlementKey string,
	entitlementName string,
	category string,
	restrictionType RestrictionType,
	tokenCost int,
	usageLimit int,
	resetPeriod *string,
) DurationAccessMapping {
	return DurationAccessMapping{
		PlanDurationID:  planDurationID,
		EntitlementID:   entitlementID,
		EntitlementKey:  entitlementKey,
		EntitlementName: entitlementName,
		Category:        category,
		RestrictionType: restrictionType,
		TokenCost:       tokenCost,
		UsageLimit:      usageLimit,
		ResetPeriod:     resetPeriod,
		Status:          MappingStatusActive,
	}
}

func (d *DurationAccessMapping) IsTokenGated() bool {
	return d.RestrictionType == RestrictionTokenGated
}

func (d *DurationAccessMapping) IsFrequencyLimited() bool {
	return d.RestrictionType == RestrictionFrequencyLimited
}
