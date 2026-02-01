package dto_response

import "time"

type ExpertIdentity struct {
	ExpertName            string `json:"expert_name"`
	ExpertProfession      string `json:"expert_profession"`
	ExpertProfessionID    int64  `json:"expert_profession_id"`
	ExpertDegree          string `json:"expert_degree"`
	ExpertExperienceYears int    `json:"expert_experience_years"`
	ExpertEducationLevel  string `json:"expert_education_level"`
	ExpertUniversity      string `json:"expert_university"`
	ExpertStudyProgram    string `json:"expert_study_program"`
}

type Top5RecommendationItem struct {
	Rank           int    `json:"rank"`
	ProfessionID   int64  `json:"profession_id"`
	ProfessionName string `json:"profession_name"`
}

type Top5RecommendationData struct {
	Status      string                   `json:"status"`
	StatusLabel string                   `json:"status_label"`
	List        []Top5RecommendationItem `json:"list"`
}

type ScoreDetail struct {
	Value int16  `json:"value"`
	Max   int    `json:"max"`
	Label string `json:"label"`
}

type ExpertScoreDetails struct {
	Accuracy   ScoreDetail `json:"accuracy"`
	Logic      ScoreDetail `json:"logic"`
	Usefulness ScoreDetail `json:"usefulness"`
}

type ExpertSuggestion struct {
	HasText bool   `json:"has_text"`
	Text    string `json:"text"`
}

type CareerProfileExpertFeedbackDetailResponse struct {
	FeedbackID          int64                  `json:"feedback_id"`
	SubmittedAt         time.Time              `json:"submitted_at"`
	TestCategory        string                 `json:"test_category"`
	TestCategoryLabel   string                 `json:"test_category_label"`
	Identity            ExpertIdentity         `json:"identity"`
	Top5Recommendations Top5RecommendationData `json:"top5_recommendations"`
	Scores              ExpertScoreDetails     `json:"scores"`
	Obstacles           []FeedbackObstacleItem `json:"obstacles"`
	Suggestion          ExpertSuggestion       `json:"suggestion"`
}
