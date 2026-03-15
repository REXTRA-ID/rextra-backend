package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementStatus string
type EntitlementLevel string
type RestrictionType string

const (
	EntitlementStatusActive   EntitlementStatus = "active"
	EntitlementStatusInactive EntitlementStatus = "inactive"

	EntitlementLevelFeature    EntitlementLevel = "fitur"
	EntitlementLevelSubFeature EntitlementLevel = "sub_fitur"

	RestrictionUnlimited        RestrictionType = "unlimited"
	RestrictionTokenGated       RestrictionType = "token_gated"
	RestrictionFrequencyLimited RestrictionType = "frequency_limited"
	RestrictionLocked           RestrictionType = "locked"
)

type Entitlement struct {
	ID               uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Key              string            `json:"key" gorm:"type:varchar(200);uniqueIndex;not null"`
	Name             string            `json:"name" gorm:"type:varchar(200);not null"`
	Description      string            `json:"description" gorm:"type:text"`
	Level            EntitlementLevel  `json:"level" gorm:"type:varchar(20);not null"`
	Status           EntitlementStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	FeatureID        uuid.UUID         `json:"feature_id" gorm:"type:uuid;not null;index"`
	SubFeatureID     *uuid.UUID        `json:"sub_feature_id,omitempty" gorm:"type:uuid;index"`
	ActionCategoryID uuid.UUID         `json:"action_category_id" gorm:"type:uuid;not null;index"`

	RestrictionType RestrictionType `json:"restriction_type" gorm:"type:varchar(30);not null;default:'unlimited'"`
	TokenCost       int             `json:"token_cost" gorm:"not null;default:0"`
	ResetPeriod     *string         `json:"reset_period,omitempty" gorm:"type:varchar(20)"`

	Timestamp

	Feature Feature `json:"feature,omitempty" gorm:"foreignKey:FeatureID;references:ID"`

	SubFeature *SubFeature `json:"sub_feature,omitempty" gorm:"foreignKey:SubFeatureID;references:ID"`

	ActionCategory ActionCategory `json:"action_category,omitempty" gorm:"foreignKey:ActionCategoryID;references:ID"`

	DurationAccessMappings []DurationAccessMapping `json:"duration_access_mappings,omitempty" gorm:"foreignKey:EntitlementID;references:ID"`
}

func (Entitlement) TableName() string {
	return "entitlements"
}

func (e *Entitlement) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now
	return nil
}

func BuildEntitlementKey(featurePrefix, subFeatureSlug, actionSlug string) string {
	if subFeatureSlug != "" {
		return featurePrefix + "." + subFeatureSlug + "." + actionSlug
	}
	return featurePrefix + "." + actionSlug
}

func NewEntitlement(
	key, name, description string,
	restrictionType, resetPeriod string,
	tokenCost int,
	level EntitlementLevel,
	featureID uuid.UUID,
	subFeatureID *uuid.UUID,
	actionCategoryID uuid.UUID,
) Entitlement {
	return Entitlement{
		Key:              key,
		Name:             name,
		Description:      description,
		Level:            level,
		FeatureID:        featureID,
		SubFeatureID:     subFeatureID,
		ActionCategoryID: actionCategoryID,
		RestrictionType:  RestrictionType(restrictionType),
		TokenCost:        tokenCost,
		ResetPeriod:      &resetPeriod,
		Status:           EntitlementStatusActive,
	}
}

func (e *Entitlement) IsTokenGated() bool {
	return e.RestrictionType == RestrictionTokenGated
}

func (e *Entitlement) IsFrequencyLimited() bool {
	return e.RestrictionType == RestrictionFrequencyLimited
}

func (e *Entitlement) IsLocked() bool {
	return e.RestrictionType == RestrictionLocked
}
