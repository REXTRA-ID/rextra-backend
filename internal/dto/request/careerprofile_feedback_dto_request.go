package dto_request

type GetStudentFeedbackRequest struct {
	TestCategory string `form:"test_category" json:"test_category"`
	Query        string `form:"q" json:"q"`
	StartDate    string `form:"start_date" json:"start_date"`
	EndDate      string `form:"end_date" json:"end_date"`
	SortBy       string `form:"sort_by" json:"sort_by"`   // submitted_at, name
	SortDir      string `form:"sort_dir" json:"sort_dir"` // asc, desc
	Page         int    `form:"page" json:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}

type GetExpertFeedbackRequest struct {
	TestCategory string `form:"test_category" json:"test_category"`
	Query        string `form:"q" json:"q"`
	Top5Status   string `form:"top5_status" json:"top5_status"`
	StartDate    string `form:"start_date" json:"start_date"`
	EndDate      string `form:"end_date" json:"end_date"`
	SortBy       string `form:"sort_by" json:"sort_by"`   // submitted_at, name
	SortDir      string `form:"sort_dir" json:"sort_dir"` // asc, desc
	Page         int    `form:"page" json:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}
