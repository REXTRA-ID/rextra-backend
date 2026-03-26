package dto_response

import (
	"time"

	"gorm.io/datatypes"
)

type CreateCareerRecommendationResponse struct {
	ID                  int64          `json:"id"`
	TestSessionID       int64          `json:"test_session_id"`
	RecommendationsData datatypes.JSON `json:"recommendations_data"`
	TopProfession1ID    *int64         `json:"top_profession_1_id,omitempty"`
	TopProfession2ID    *int64         `json:"top_profession_2_id,omitempty"`
	GeneratedAt         time.Time      `json:"generated_at"`
	AIModelUsed         string         `json:"ai_model_used"`
}
