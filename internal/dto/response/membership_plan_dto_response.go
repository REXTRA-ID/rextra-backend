package dto_response

type GetMembershipPlanResponse struct {
	ID           string  `json:"id"`
	PlanName     string  `json:"plan_name"`
	Category     string  `json:"category"`
	TierLabel    string  `json:"tier_label"`
	EmblemKey    string  `json:"emblem_key"`
	Description  string  `json:"description"`
	Status       string  `json:"status"`
	PricingMode  string  `json:"pricing_mode"`
	DurationMode string  `json:"duration_mode"`
	BasePrice1M  int64   `json:"base_price_1m"`
	BaseToken1M  int     `json:"base_token_1m"`
	Discount3M   float64 `json:"discount_3m"`
	Discount6M   float64 `json:"discount_6m"`
	Discount12M  float64 `json:"discount_12m"`
	BonusToken3M  int    `json:"bonus_token_3m"`
	BonusToken6M  int    `json:"bonus_token_6m"`
	BonusToken12M int    `json:"bonus_token_12m"`
	ActiveUsers  int     `json:"active_users"`
	Benefits     []string `json:"benefits"`
	StarterDurationMonths int `json:"starter_duration_months"`
	PlanDurations []GetPlanDurationResponse `json:"plan_durations,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetPlanDurationResponse struct {
	ID             string  `json:"id"`
	PlanID         string  `json:"plan_id"`
	DurationMonths int     `json:"duration_months"`
	Price          int64   `json:"price"`
	DiscountPct    float64 `json:"discount_pct"`
	FinalPrice     int64   `json:"final_price"`
	DurationPrice  int64   `json:"duration_price"`
	TokenAmount    int     `json:"token_amount"`
	BonusToken     int     `json:"bonus_token"`
	PointsActive   bool    `json:"points_active"`
	PointsValue    int     `json:"points_value"`
	BonusPoints    int     `json:"bonus_points"`
	IsActive       bool    `json:"is_active"`
	MappingCount   int     `json:"mapping_count"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type GetDurationAccessMappingResponse struct {
	ID              string `json:"id"`
	PlanDurationID  string `json:"plan_duration_id"`
	EntitlementID   string `json:"entitlement_id"`
	EntitlementKey  string `json:"entitlement_key"`
	EntitlementName string `json:"entitlement_name"`
	Category        string `json:"category"`
	RestrictionType string `json:"restriction_type"`
	TokenCost       int    `json:"token_cost"`
	UsageLimit      int    `json:"usage_limit"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
}
