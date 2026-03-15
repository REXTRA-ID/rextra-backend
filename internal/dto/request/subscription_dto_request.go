package dto_request

type UserMembershipFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	PlanName string `form:"plan_name" binding:"omitempty,oneof=Standard Starter Basic Pro Max"`
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type SubscriptionCycleFilterRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=0"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	PlanName  string `form:"plan_name" binding:"omitempty"`
	Status    string `form:"status" binding:"omitempty,oneof=active completed expired"`
	DateFrom  string `form:"date_from" binding:"omitempty"`
	DateTo    string `form:"date_to" binding:"omitempty"`
	UserID    string `form:"user_id" binding:"omitempty"`
}
