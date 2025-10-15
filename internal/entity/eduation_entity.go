package entity

import "github.com/google/uuid"

type EducationLevel string

const (
	D1 EducationLevel = "D1"
	D2 EducationLevel = "D2"
	D3 EducationLevel = "D3"
	D4 EducationLevel = "D4"
	S1 EducationLevel = "S1"
	S2 EducationLevel = "S2"
	S3 EducationLevel = "S3"
)

type EducationStatus string

const (
	ACTIVE EducationStatus = "ACTIVE"
	GRADUATED EducationStatus = "GRADUATED"
	DROPPED EducationStatus = "DROPPED"
	DEFERRED EducationStatus = "DEFERRED"
	TRANSFERRED EducationStatus = "TRANSFERRED"
	DISMISSED EducationStatus = "DISMISSED"
)

type Education struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	
	InstitutionName string `json:"institution_name" gorm:"not null"`
	Major string `json:"major" gorm:"not null"`
	Faculty string `json:"faculty" gorm:"not null"`
	EntryYear int `json:"entry_year" gorm:"not null"`

	ExpectedGraduationYear int `json:"expected_graduation_year" gorm:"not null"`
	ActualGraduationYear int `json:"actual_graduation_year" gorm:"default:null"`

	CurrentSemester int `json:"current_semester" gorm:"not null"`
	TotalSemester int `json:"total_semester" gorm:"not null"`

	EducationLevel EducationLevel `json:"education_level" gorm:"not null"`
	Status EducationStatus `json:"status" gorm:"not null"`
	IsActive bool `json:"is_active" gorm:"default:false;not null"`

	Timestamp
}
