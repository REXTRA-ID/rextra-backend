package dto_response

type (
	FeedbackListResponse struct {
		Data       []FeedbackItem `json:"data"`
		Pagination PaginationMeta `json:"pagination"`
	}

	FeedbackItem struct {
		ID                int64    `json:"id"`
		UserName          string   `json:"user_name"`
		EaseOfUseScore    int      `json:"ease_of_use_score"`
		RelevanceScore    int      `json:"relevance_score"`
		SatisfactionScore int      `json:"satisfaction_score"`
		Obstacles         []string `json:"obstacles"`
		SubmittedAt       string   `json:"submitted_at"`
	}

	FeedbackStatsResponse struct {
		TotalFeedback     int                    `json:"total_feedback"`
		AvgEaseOfUse      float64                `json:"avg_ease_of_use"`
		AvgRelevance      float64                `json:"avg_relevance"`
		AvgSatisfaction   float64                `json:"avg_satisfaction"`
		ParticipationRate float64                `json:"participation_rate"`
		TrendData         map[string]interface{} `json:"trend_data"`
	}

	ExpertFeedbackListResponse struct {
		Data       []ExpertFeedbackItem `json:"data"`
		Pagination PaginationMeta       `json:"pagination"`
	}

	ExpertFeedbackItem struct {
		ID             int64    `json:"id"`
		ExpertName     string   `json:"expert_name"`
		Profession     string   `json:"profession"`
		TopNStatus     string   `json:"top_n_status"`
		AccuracyScore  int      `json:"accuracy_score"`
		LogicScore     int      `json:"logic_score"`
		BenefitScore   int      `json:"benefit_score"`
		Obstacles      []string `json:"obstacles"`
		SubmittedAt    string   `json:"submitted_at"`
		Recommendation *string  `json:"recommendation,omitempty"`
	}

	ExpertFeedbackDetailResponse struct {
		ID                  int64    `json:"id"`
		ExpertName          string   `json:"expert_name"`
		Profession          string   `json:"profession"`
		Degree              string   `json:"degree"`
		Experience          string   `json:"experience"`
		Education           string   `json:"education"`
		University          string   `json:"university"`
		StudyProgram        string   `json:"study_program"`
		CategoryTest        string   `json:"category_test"`
		TopFiveProfessions  []string `json:"top_five_professions"`
		TopNStatus          string   `json:"top_n_status"`
		AccuracyScore       int      `json:"accuracy_score"`
		LogicScore          int      `json:"logic_score"`
		BenefitScore        int      `json:"benefit_score"`
		Obstacles           []string `json:"obstacles"`
		Suggestions         string   `json:"suggestions"`
		SubmittedAt         string   `json:"submitted_at"`
		TestSessionToken    string   `json:"test_session_token,omitempty"`
		RecommendationNotes string   `json:"recommendation_notes,omitempty"`
	}
)
