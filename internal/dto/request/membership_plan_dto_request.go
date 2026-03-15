package dto_request

type CreateMembershipPlanRequest struct {
	PlanName     string  `json:"plan_name" binding:"required,oneof=Standard Starter Basic Pro Max"`
	Category     string  `json:"category" binding:"required,oneof=unpaid paid"`
	TierLabel    string  `json:"tier_label" binding:"required"`
	EmblemKey    string  `json:"emblem_key"`
	Description  string  `json:"description"`
	MarketingIntro string `json:"marketing_intro"`
	PricingMode  string  `json:"pricing_mode" binding:"required,oneof=manual otomatis"`
	DurationMode string  `json:"duration_mode" binding:"required,oneof=dengan_durasi tanpa_durasi"`
	BasePrice1M  int64   `json:"base_price_1m"`
	BaseToken1M  int     `json:"base_token_1m"`
	Discount3M   float64 `json:"discount_3m"`
	Discount6M   float64 `json:"discount_6m"`
	Discount12M  float64 `json:"discount_12m"`
	BonusToken3M  int    `json:"bonus_token_3m"`
	BonusToken6M  int    `json:"bonus_token_6m"`
	BonusToken12M int    `json:"bonus_token_12m"`
	Benefits     []string `json:"benefits"`
}

type UpdateMembershipPlanRequest struct {
	TierLabel    string   `json:"tier_label" binding:"required"`
	EmblemKey    string   `json:"emblem_key"`
	Description  string   `json:"description"`
	MarketingIntro string `json:"marketing_intro"`
	Status       string   `json:"status" binding:"required,oneof=aktif nonaktif"`
	PricingMode  string   `json:"pricing_mode" binding:"required,oneof=manual otomatis"`
	BasePrice1M  int64    `json:"base_price_1m"`
	BaseToken1M  int      `json:"base_token_1m"`
	Discount3M   float64  `json:"discount_3m"`
	Discount6M   float64  `json:"discount_6m"`
	Discount12M  float64  `json:"discount_12m"`
	BonusToken3M  int     `json:"bonus_token_3m"`
	BonusToken6M  int     `json:"bonus_token_6m"`
	BonusToken12M int     `json:"bonus_token_12m"`
	Benefits     []string `json:"benefits"`
	StarterDurationMonths int `json:"starter_duration_months" binding:"min=0"`
}
