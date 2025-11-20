package dto_response

type GetMembershipPlanResponse struct {
	ID               string  `json:"id"`
	PlanName         string  `json:"plan_name"`
	MonthlyToken     int     `json:"monthly_token"`
	BaseMonthlyPrice float64 `json:"base_monthly_price"`
	Description      string  `json:"description"`
	Benefits         []byte  `json:"benefits"`
	IsActive         bool    `json:"is_active"`
}
