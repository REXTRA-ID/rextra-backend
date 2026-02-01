package entity

import "gorm.io/datatypes"

type CareerProfileFeedbackExpert struct {
	FeedbackID      int64 `json:"feedback_id" gorm:"primaryKey"`
	AccuracyScore   int16 `json:"accuracy_score" gorm:"not null"`   // 1-7 scale
	LogicScore      int16 `json:"logic_score" gorm:"not null"`      // 1-7 scale
	UsefulnessScore int16 `json:"usefulness_score" gorm:"not null"` // 1-7 scale

	// Expert identity snapshot fields
	ExpertName            string `json:"expert_name" gorm:"size:255;not null"`
	ExpertProfession      string `json:"expert_profession" gorm:"size:255"`
	ExpertDegree          string `json:"expert_degree" gorm:"size:255"`
	ExpertExperienceYears *int32 `json:"expert_experience_years"`
	ExpertEducationLevel  string `json:"expert_education_level" gorm:"size:255"`
	ExpertUniversity      string `json:"expert_university" gorm:"size:255"`
	ExpertStudyProgram    string `json:"expert_study_program" gorm:"size:255"`

	// Recommendation validation
	ExpertProfessionID      *int64         `json:"expert_profession_id"` // ID profesi yang dipilih expert dari rekomendasi
	Top5RecommendationsJSON datatypes.JSON `json:"top5_recommendations_json" gorm:"type:jsonb;not null;default:'[]'"`
	Top5Status              string         `json:"top5_status" gorm:"type:varchar(20);->;virtualType:GENERATED ALWAYS AS (CASE WHEN expert_profession_id IS NULL THEN 'NOT_PRESENT' WHEN top5_recommendations_json->0->>'profession_id' = expert_profession_id::text THEN 'P1' WHEN top5_recommendations_json->1->>'profession_id' = expert_profession_id::text THEN 'P2' WHEN top5_recommendations_json->2->>'profession_id' = expert_profession_id::text OR top5_recommendations_json->3->>'profession_id' = expert_profession_id::text OR top5_recommendations_json->4->>'profession_id' = expert_profession_id::text THEN 'P3_5' ELSE 'NOT_PRESENT' END) STORED"`

	SuggestionText *string `json:"suggestion_text" gorm:"type:text"`

	Feedback  KenaliDiriFeedback                    `json:"-" gorm:"foreignKey:FeedbackID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Obstacles []CareerProfileFeedbackExpertObstacle `json:"obstacles" gorm:"foreignKey:FeedbackID;references:FeedbackID"`
}

func (CareerProfileFeedbackExpert) TableName() string {
	return "careerprofile_feedback_expert"
}
