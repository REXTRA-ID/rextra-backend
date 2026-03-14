package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActionCategoryStatus string

const (
	ActionCategoryStatusActive   ActionCategoryStatus = "active"
	ActionCategoryStatusInactive ActionCategoryStatus = "inactive"
)

type ActionCategory struct {
	ID          uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string               `json:"name" gorm:"type:varchar(50);not null"`
	Slug        string               `json:"slug" gorm:"type:varchar(50);uniqueIndex;not null"`
	Description string               `json:"description" gorm:"type:text"`
	Status      ActionCategoryStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`

	Timestamp

	// HasMany Entitlement — action category dipakai oleh banyak entitlement sebagai komponen terakhir key.
	Entitlements []Entitlement `json:"entitlements,omitempty" gorm:"foreignKey:ActionCategoryID;references:ID"`
}

func (ActionCategory) TableName() string {
	return "action_categories"
}

func (a *ActionCategory) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	return nil
}

func NewActionCategory(name, slug, description, status string) ActionCategory {
	return ActionCategory{
		Name:        name,
		Slug:        slug,
		Description: description,
		Status:      ActionCategoryStatus(status),
	}
}
