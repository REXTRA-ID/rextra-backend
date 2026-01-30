package dto_request

// GetTestHistoryRequest carries filters and pagination for admin test history list.
type GetTestHistoryRequest struct {
	Page       int     `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit      int     `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	CategoryID *int64  `form:"category_id" json:"category_id"`
	Status     string  `form:"status" json:"status" binding:"omitempty,oneof=ongoing completed abandoned"`
	UserName   string  `form:"user_name" json:"user_name"`
	StartDate  string  `form:"start_date" json:"start_date"`
	EndDate    string  `form:"end_date" json:"end_date"`
	SortBy     string  `form:"sort_by" json:"sort_by" binding:"omitempty,oneof=name_asc name_desc date_asc date_desc"`
}

// ExportTestHistoryRequest wraps filters for exporting test history.
type ExportTestHistoryRequest struct {
	CategoryID *int64 `form:"category_id" json:"category_id"`
	Status     string `form:"status" json:"status"`
	StartDate  string `form:"start_date" json:"start_date"`
	EndDate    string `form:"end_date" json:"end_date"`
	Format     string `form:"format" json:"format" binding:"required,oneof=csv excel pdf"`
}

// DeleteTestDataRequest supports single or bulk deletion by test IDs.
type DeleteTestDataRequest struct {
	TestIDs []int64 `json:"test_ids" binding:"required,min=1,dive,required"`
}

// UpdateRiasecCodeRequest carries payload to update master data.
type UpdateRiasecCodeRequest struct {
	RiasecTitle       string   `json:"riasec_title" binding:"required"`
	RiasecDescription string   `json:"riasec_description" binding:"required"`
	Strengths         []string `json:"strengths" binding:"required,min=1,dive,required"`
	Challenges        []string `json:"challenges" binding:"required,min=1,dive,required"`
	Strategies        []string `json:"strategies" binding:"required,min=1,dive,required"`
	WorkEnvironments  []string `json:"work_environments" binding:"required,min=1,dive,required"`
	InteractionStyles []string `json:"interaction_styles" binding:"required,min=1,dive,required"`
}

// GetRiasecCodeListRequest filters RIASEC master data.
type GetRiasecCodeListRequest struct {
	CodeType string `form:"code_type" json:"code_type" binding:"omitempty,oneof=single dual triple"`
	Search   string `form:"search" json:"search"`
}
