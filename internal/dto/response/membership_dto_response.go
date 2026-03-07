package dto_response

type GetMembershipPlanResponse struct {
	ID               string   `json:"id"`
	PlanName         string   `json:"plan_name"`
	MonthlyToken     int      `json:"monthly_token"`
	BaseMonthlyPrice float64  `json:"base_monthly_price"`
	Description      string   `json:"description"`
	Benefits         []string `json:"benefits"`
	IsActive         bool     `json:"is_active"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetMembershipDurationResponse struct {
	ID                   string  `json:"id"`
	DurationMonth        int     `json:"duration_month"`
	TokenBonusPercentage float64 `json:"token_bonus_percentage"`
	RextraPoinMultiplier int     `json:"rextra_poin_multiplier"`
	IsActive             bool    `json:"is_active"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetDurationAccessMappingResponse struct {
	ID              string `json:"id"`
	PlanDurationID  string `json:"plan_duration_id"`
	EntitlementID   string `json:"entitlement_id"`
	EntitlementKey  string `json:"entitlement_key"`
	EntitlementName string `json:"entitlement_name"`
	Category        string `json:"category"`

	RestrictionType string  `json:"restriction_type"`
	TokenCost       int     `json:"token_cost"`
	UsageLimit      int     `json:"usage_limit"`
	ResetPeriod     *string `json:"reset_period,omitempty"`

	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
