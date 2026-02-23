package entity

import "github.com/google/uuid"

type FeatureStatus string
type FeatureType string

const (
	FEATURE_AKTIF    FeatureStatus = "Aktif"
	FEATURE_NONAKTIF FeatureStatus = "Nonaktif"

	MAINFEATURE FeatureType = "Fitur"
	SUBFEATURE  FeatureType = "Sub Fitur"
)

type Feature struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	Name   string        `json:"name"`
	Slug   string        `json:"slug"`
	Status FeatureStatus `json:"status"`
	Type   FeatureType   `json:"type"`

	ParentID *uuid.UUID `json:"parent_id"`
	Parent   *Feature   `gorm:"foreignKey:ParentID"`
}

func NewFeature(name, slug string, parentId *uuid.UUID, status string) Feature {
	var fType FeatureType = MAINFEATURE
	if parentId != nil {
		fType = SUBFEATURE
	}

	return Feature{
		Name:     name,
		Slug:     slug,
		ParentID: parentId,
		Status:   FeatureStatus(status),
		Type:     fType,
	}
}
