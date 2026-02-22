package entity

import "github.com/google/uuid"

type ActionStatus string

const (
	ACTION_AKTIF    ActionStatus = "Aktif"
	ACTION_NONAKTIF ActionStatus = "Nonaktif"
)

type ActionCategory struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Description string       `json:"description"`
	Status      ActionStatus `json:"status"`
}

func NewActionCategory(name string, slug string, description string, status string) ActionCategory {
	return ActionCategory{
		Name:        name,
		Slug:        slug,
		Description: description,
		Status:      ActionStatus(status),
	}
}
