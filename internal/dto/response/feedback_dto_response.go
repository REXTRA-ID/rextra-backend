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
		TotalFeedback        int                        `json:"total_feedback"`
		AvgEaseOfUse         float64                    `json:"avg_ease_of_use"`
		AvgRelevance         float64                    `json:"avg_relevance"`
		AvgSatisfaction      float64                    `json:"avg_satisfaction"`
		ParticipationRate    float64                    `json:"participation_rate"`
		TrendData            TrendChartData             `json:"trend_data"`
		ScoreDistribution    ScoreDistributionSet       `json:"score_distribution"`
		ObstacleSummary      ObstacleSummary            `json:"obstacle_summary"`
		ObstacleDistribution []ObstacleDistributionItem `json:"obstacle_distribution"`
		SentimentComposition []SentimentCompositionItem `json:"sentiment_composition"`
		ResponseRateInfo     ResponseRateInfo           `json:"response_rate_info"`
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

	TrendChartData struct {
		Labels         []string `json:"labels"`
		TestCounts     []int    `json:"test_counts"`
		FeedbackCounts []int    `json:"feedback_counts"`
	}

	ScoreDistributionItem struct {
		Score      int     `json:"score"`
		Count      int     `json:"count"`
		Percentage float64 `json:"percentage"`
	}

	ScoreDistributionSet struct {
		EaseOfUse    []ScoreDistributionItem `json:"ease_of_use"`
		Relevance    []ScoreDistributionItem `json:"relevance"`
		Satisfaction []ScoreDistributionItem `json:"satisfaction"`
	}

	ObstacleSummary struct {
		WithObstacles    int `json:"with_obstacles"`
		WithoutObstacles int `json:"without_obstacles"`
	}

	ObstacleDistributionItem struct {
		Name       string  `json:"name"`
		Count      int     `json:"count"`
		Percentage float64 `json:"percentage"`
	}

	SentimentCompositionItem struct {
		Metric   string  `json:"metric"`
		Negative float64 `json:"negative"`
		Neutral  float64 `json:"neutral"`
		Positive float64 `json:"positive"`
	}

	ResponseRateInfo struct {
		FeedbackCount  int     `json:"feedback_count"`
		CompletedTests int     `json:"completed_tests"`
		NotFilled      int     `json:"not_filled"`
		ResponseRate   float64 `json:"response_rate"`
	}
)
