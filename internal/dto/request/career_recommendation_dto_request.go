package dto_request

type (
	CreateCareerRecommendationRequest struct {
		UserID         string                                `json:"user_id"`
		Analysis       []CareerRecommendationAnalysisRequest `json:"analysis" binding:"required"`
		TopProfessions []string                              `json:"top_2" binding:"required"`
	}

	CareerRecommendationAnalysisRequest struct {
		Profession      string `json:"profession" binding:"required"`
		MatchPercentage int    `json:"match_percentage" binding:"required"`
		Reason          string `json:"reason" binding:"required"`
	}
)
