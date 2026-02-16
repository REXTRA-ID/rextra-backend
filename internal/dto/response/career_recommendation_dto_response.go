package dto_response

type (
	CreateCareerRecommendationResponse struct {
		ID             string                                 `json:"id"`
		UserID         string                                 `json:"user_id"`
		Analysis       []CareerRecommendationAnalysisResponse `json:"analysis"`
		TopProfessions []string                               `json:"top_2"`
	}

	CareerRecommendationAnalysisResponse struct {
		Profession      string `json:"profession"`
		MatchPercentage int    `json:"match_percentage"`
		Reason          string `json:"reason"`
	}
)
