package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeatureStatus string
type FeatureType string

const (
	FeatureStatusActive   FeatureStatus = "active"
	FeatureStatusInactive FeatureStatus = "inactive"

	FeatureTypeSingle     FeatureType = "tunggal"    // fitur tanpa sub fitur
	FeatureTypeHierarchic FeatureType = "bertingkat" // fitur dengan sub fitur
)

type Feature struct {
	ID          uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string        `json:"name" gorm:"type:varchar(100);not null"`
	Slug        string        `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	Prefix      string        `json:"prefix" gorm:"type:varchar(20);uniqueIndex;not null"`
	Description string        `json:"description" gorm:"type:text"`
	Type        FeatureType   `json:"type" gorm:"type:varchar(20);not null;default:'tunggal'"`
	Status      FeatureStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	Timestamp

	// HasMany SubFeature — fitur bertingkat memiliki banyak sub fitur di bawahnya.
	SubFeatures []SubFeature `json:"sub_features,omitempty" gorm:"foreignKey:FeatureID;references:ID"`

	// HasMany Entitlement — setiap entitlement punya satu fitur induk sebagai "pemilik" prefix-nya.
	Entitlements []Entitlement `json:"entitlements,omitempty" gorm:"foreignKey:FeatureID;references:ID"`
}

func (Feature) TableName() string {
	return "features"
}

func (f *Feature) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	f.CreatedAt = now
	f.UpdatedAt = now
	return nil
}

func NewFeature(name, slug, prefix, description string, featureType FeatureType) Feature {
	return Feature{
		Name:        name,
		Slug:        slug,
		Prefix:      prefix,
		Description: description,
		Type:        featureType,
		Status:      FeatureStatusActive,
	}
}
