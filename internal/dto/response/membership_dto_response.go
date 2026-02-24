package dto_response

type GetMembershipPlanResponse struct {
	ID               string   `json:"id"`
	PlanName         string   `json:"plan_name"`
	MonthlyToken     int      `json:"monthly_token"`
	BaseMonthlyPrice float64  `json:"base_monthly_price"`
	Description      string   `json:"description"`
	Benefits         []string `json:"benefits"`
	IsActive         bool     `json:"is_active"`
}

type GetMembershipDurationResponse struct {
	ID                   string  `json:"id"`
	DurationMonth        int     `json:"duration_month"`
	TokenBonusPercentage float64 `json:"token_bonus"`
	RextraPoinMultiplier int     `json:"rextra_point_multiplier"`
	IsActive             bool    `json:"is_active"`
}
