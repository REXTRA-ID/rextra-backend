package dto_response

type (
	TestHistoryListResponse struct {
		Data       []TestHistoryItem `json:"data"`
		Pagination PaginationMeta    `json:"pagination"`
	}

	TestHistoryItem struct {
		TestID         string  `json:"test_id"`
		UserName       string  `json:"user_name"`
		CategoryName   string  `json:"category_name"`
		Status         string  `json:"status"`
		ResultCode     string  `json:"result_code"`
		StartedAt      string  `json:"started_at"`
		CompletedAt    *string `json:"completed_at,omitempty"`
		RiasecCodeType string  `json:"riasec_code_type,omitempty"`
	}

	TestDetailResponse struct {
		TestID          string                 `json:"test_id"`
		UserName        string                 `json:"user_name"`
		CategoryName    string                 `json:"category_name"`
		Status          string                 `json:"status"`
		StartedAt       string                 `json:"started_at"`
		CompletedAt     *string                `json:"completed_at,omitempty"`
		RiasecResult    *RiasecResultDetail    `json:"riasec_result,omitempty"`
		IkigaiResult    *IkigaiResultDetail    `json:"ikigai_result,omitempty"`
		Recommendations []RecommendationDetail `json:"recommendations,omitempty"`
	}

	RiasecResultDetail struct {
		ScoreR            int    `json:"score_r"`
		ScoreI            int    `json:"score_i"`
		ScoreA            int    `json:"score_a"`
		ScoreS            int    `json:"score_s"`
		ScoreE            int    `json:"score_e"`
		ScoreC            int    `json:"score_c"`
		RiasecCode        string `json:"riasec_code"`
		RiasecTitle       string `json:"riasec_title"`
		ClassificationType string `json:"classification_type"`
		IsInconsistent    bool   `json:"is_inconsistent"`
	}

	IkigaiResultDetail struct {
		LoveNarrative       string `json:"love_narrative"`
		GoodAtNarrative     string `json:"good_at_narrative"`
		WorldNeedsNarrative string `json:"world_needs_narrative"`
		PaidForNarrative    string `json:"paid_for_narrative"`
	}

	RecommendationDetail struct {
		Rank            int    `json:"rank"`
		ProfessionID    int64  `json:"profession_id"`
		ProfessionName  string `json:"profession_name"`
		MatchPercentage int    `json:"match_percentage"`
		MatchReasoning  string `json:"match_reasoning"`
	}

	RiasecCodeListResponse struct {
		Data       []RiasecCodeItem `json:"data"`
		TotalCodes int              `json:"total_codes"`
	}

	RiasecCodeItem struct {
		ID       int64  `json:"id"`
		Code     string `json:"code"`
		Title    string `json:"title"`
		CodeType string `json:"code_type"`
	}

	RiasecCodeDetailResponse struct {
		ID                int64    `json:"id"`
		Code              string   `json:"code"`
		Title             string   `json:"title"`
		Description       string   `json:"description"`
		Strengths         []string `json:"strengths"`
		Challenges        []string `json:"challenges"`
		Strategies        []string `json:"strategies"`
		WorkEnvironments  []string `json:"work_environments"`
		InteractionStyles []string `json:"interaction_styles"`
	}

	ExportFileResponse struct {
		FileName     string `json:"file_name"`
		FileURL      string `json:"file_url"`
		ExportedAt   string `json:"exported_at"`
		TotalRecords int    `json:"total_records"`
	}
)
