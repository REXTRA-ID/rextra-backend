package dto_response

import "time"

type FeedbackObstacleItem struct {
	Key       string  `json:"key"`
	Label     string  `json:"label"`
	OtherText *string `json:"other_text"`
}

type CareerProfileStudentItem struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type CareerProfileStudentScores struct {
	Ease         int16 `json:"ease"`
	Relevance    int16 `json:"relevance"`
	Satisfaction int16 `json:"satisfaction"`
}

type CareerProfileStudentFeedbackItem struct {
	ID                int64                      `json:"id"`
	Student           CareerProfileStudentItem   `json:"student"`
	Scores            CareerProfileStudentScores `json:"scores"`
	Obstacles         []FeedbackObstacleItem     `json:"obstacles"`
	MessageToTeam     *string                    `json:"message_to_team"`
	SubmittedAt       time.Time                  `json:"submitted_at"`
	TestCategory      string                     `json:"test_category"`
	TestCategoryLabel string                     `json:"test_category_label"`
}

type CareerProfileExpertItem struct {
	Name         string `json:"name"`
	Profession   string `json:"profession"`
	ProfessionID int64  `json:"profession_id"`
}

type CareerProfileExpertScores struct {
	Accuracy   int16 `json:"accuracy"`
	Logic      int16 `json:"logic"`
	Usefulness int16 `json:"usefulness"`
}

type CareerProfileExpertFeedbackItem struct {
	ID                int64                     `json:"id"`
	Expert            CareerProfileExpertItem   `json:"expert"`
	Top5Status        string                    `json:"top5_status"`
	Top5StatusLabel   string                    `json:"top5_status_label"`
	Scores            CareerProfileExpertScores `json:"scores"`
	Obstacles         []FeedbackObstacleItem    `json:"obstacles"`
	HasSuggestion     bool                      `json:"has_suggestion"`
	SubmittedAt       time.Time                 `json:"submitted_at"`
	TestCategory      string                    `json:"test_category"`
	TestCategoryLabel string                    `json:"test_category_label"`
}

type CareerProfilePaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type CareerProfileStudentFeedbackListResponse struct {
	Items      []CareerProfileStudentFeedbackItem `json:"items"`
	Pagination CareerProfilePaginationMeta        `json:"pagination"`
}

type CareerProfileExpertFeedbackListResponse struct {
	Items      []CareerProfileExpertFeedbackItem `json:"items"`
	Pagination CareerProfilePaginationMeta       `json:"pagination"`
}
type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type SortOption struct {
	Value   string `json:"value"`
	Label   string `json:"label"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

type CareerProfileFeedbackMetadataResponse struct {
	TestCategories []LabelValue `json:"test_categories"`
	SortOptions    []SortOption `json:"sort_options"`
	Top5Statuses   []LabelValue `json:"top5_statuses"`
}
