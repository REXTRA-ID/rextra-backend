package dto_request

type GetFeedbackListRequest struct {
	Page         int     `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit        int     `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	CategoryID   *int64  `form:"category_id" json:"category_id"`
	UserName     string  `form:"user_name" json:"user_name"`
	HasObstacles *bool   `form:"has_obstacles" json:"has_obstacles"`
	SortBy       string  `form:"sort_by" json:"sort_by"`
}

type GetExpertFeedbackListRequest struct {
	Page       int     `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit      int     `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	CategoryID *int64  `form:"category_id" json:"category_id"`
	ExpertName string  `form:"expert_name" json:"expert_name"`
	TopNStatus string  `form:"top_n_status" json:"top_n_status" binding:"omitempty,oneof=P1 P2 P3-5 not_found"`
	SortBy     string  `form:"sort_by" json:"sort_by"`
}
