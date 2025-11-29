package entity

import "github.com/google/uuid"

type PersonaType string

const (
	Pathfinder PersonaType = "pathfinder"
	Builder    PersonaType = "builder"
	Achiever   PersonaType = "achiever"
)

type Persona struct {
	ID                        uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID                    uuid.UUID   `json:"user_id" gorm:"not null"`
	PersonaType               PersonaType `json:"persona_type" gorm:"default:'pathfinder';not null"`
	EducationSaved            bool        `json:"education_saved" gorm:"default:false;not null"`
	CareerRecommendationTired bool        `json:"career_recommendation_tired" gorm:"default:false;not null"`
	CareerDictionaryAccessed  bool        `json:"career_dictionary_accessed" gorm:"default:false;not null"`
	CareerPlanCreated         bool        `json:"career_plan_created" gorm:"default:false;not null"`
	PorfolioRecorded          bool        `json:"portfolio_recorded" gorm:"default:false;not null"`
	ExplorationAIUsed         bool        `json:"exploration_ai_used" gorm:"default:false;not null"`
	CVCreated                 bool        `json:"cv_created" gorm:"default:false;not null"`
	InterviewSimulated        bool        `json:"interview_simulated" gorm:"default:false;not null"`
	LinkedinOptimaze          bool        `json:"linkedin_optimaze" gorm:"default:false;not null"`
	IntershipPlanReported     bool        `json:"intership_plan_reported" gorm:"default:false;not null"`
	Timestamp
}

func (p *Persona) TableName() string {
	return "personas"
}
