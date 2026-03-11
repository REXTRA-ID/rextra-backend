package dto_request

import "gorm.io/datatypes"

type CreateCareerRecommendationRequest struct {
	TestSessionID       int64          `json:"test_session_id" binding:"required"`
	RecommendationsData datatypes.JSON `json:"recommendations_data" binding:"required"`
	TopProfession1ID    *int64         `json:"top_profession_1_id"`
	TopProfession2ID    *int64         `json:"top_profession_2_id"`
	AIModelUsed         string         `json:"ai_model_used"`
}
