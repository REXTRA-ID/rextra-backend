package entity

import "github.com/google/uuid"

type PersonaType string

const (
	Pathfinder PersonaType = "Pathfinder"
)

type Persona struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `json:"user_id" gorm:"not null"`
	Institution string    `json:"institution" gorm:"not null"`
	Study       string    `json:"study" gorm:"not null"`

	// enumerated [0 = D3, 1= D4/S1, 2 = S2, 3 = S3]
	EducationLevel string `json:"education_level" gorm:"not null"`
	GraduationYear int    `json:"graduation_year" gorm:"not null"`
	CareerPlan     string `json:"career_plan" gorm:"not null"`
	CareerDreams   string `json:"career_dreams" gorm:"not null"`
	Portfolio      bool   `json:"portfolio" gorm:"default:false;not null"`
	Application    bool   `json:"application" gorm:"not null"`

	// Additional fields needed by seeder
	PersonaType            PersonaType `json:"persona_type" gorm:"size:100"`
	EducationSaved         bool        `json:"education_saved" gorm:"default:false"`
	CareerRecommendationTired bool     `json:"career_recommendation_tired" gorm:"default:false"`
	CareerDictionaryAccessed bool      `json:"career_dictionary_accessed" gorm:"default:false"`
	CareerPlanCreated      bool        `json:"career_plan_created" gorm:"default:false"`
	PorfolioRecorded       bool        `json:"porfolio_recorded" gorm:"default:false"`
	ExplorationAIUsed      bool        `json:"exploration_ai_used" gorm:"default:false"`
	CVCreated              bool        `json:"cv_created" gorm:"default:false"`
	InterviewSimulated     bool        `json:"interview_simulated" gorm:"default:false"`
	LinkedinOptimaze       bool        `json:"linkedin_optimaze" gorm:"default:false"`
	IntershipPlanReported  bool        `json:"intership_plan_reported" gorm:"default:false"`

	// enumerated [0 = mahasiswa aktif, 1= fresh graduate, 2 = professional]
	Status string `json:"status" gorm:"not null"`

	Timestamp
}

func (p *Persona) TableName() string {
	return "personas"
}
