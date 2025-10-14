package entity

import "github.com/google/uuid"

type EducationPlan struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	
	InstitutionName string `json:"institution_name" gorm:"not null"`
	Major string `json:"major" gorm:"not null"`
	Faculty string `json:"faculty" gorm:"not null"`
	ExpectedEntryYear int `json:"expected_entry_year" gorm:"not null"`
	EducationProgram string `json:"education_program" gorm:"not null"`
	
	EducationLevel EducationLevel `json:"education_level" gorm:"not null"`

	Timestamp
}
