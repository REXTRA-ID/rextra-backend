package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubFeatureStatus string

const (
	SubFeatureStatusActive   SubFeatureStatus = "active"
	SubFeatureStatusInactive SubFeatureStatus = "inactive"
)

type SubFeature struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FeatureID   uuid.UUID        `json:"feature_id" gorm:"type:uuid;not null;index"`
	Name        string           `json:"name" gorm:"type:varchar(100);not null"`
	Slug        string           `json:"slug" gorm:"type:varchar(100);not null"`
	Description string           `json:"description" gorm:"type:text"`
	Status      SubFeatureStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	Timestamp

	Feature      Feature       `json:"feature,omitempty" gorm:"foreignKey:FeatureID;references:ID"`
	Entitlements []Entitlement `json:"entitlements,omitempty" gorm:"foreignKey:SubFeatureID;references:ID"`
}

func (SubFeature) TableName() string { return "sub_features" }

func (s *SubFeature) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

func NewSubFeature(featureID uuid.UUID, name, slug, description string) SubFeature {
	return SubFeature{
		FeatureID:   featureID,
		Name:        name,
		Slug:        slug,
		Description: description,
		Status:      SubFeatureStatusActive,
	}
}
