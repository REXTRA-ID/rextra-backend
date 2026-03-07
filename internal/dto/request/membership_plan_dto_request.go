package dto_request

type CreateMembershipPlanRequest struct {
	PlanName         string   `json:"plan_name" binding:"required,oneof=Starter Basic Pro Max Standard"`
	MonthlyToken     int      `json:"monthly_token" binding:"required,min=0"`
	BaseMonthlyPrice float64  `json:"base_monthly_price" binding:"required,min=0"`
	Description      string   `json:"description"`
	Benefits         []string `json:"benefits"`
}

type UpdateMembershipPlanRequest struct {
	PlanName         string   `json:"plan_name" binding:"required,oneof=Starter Basic Pro Max Standard"`
	MonthlyToken     int      `json:"monthly_token" binding:"required,min=0"`
	BaseMonthlyPrice float64  `json:"base_monthly_price" binding:"required,min=0"`
	Description      string   `json:"description"`
	Benefits         []string `json:"benefits"`
	IsActive         bool     `json:"is_active"`
}
