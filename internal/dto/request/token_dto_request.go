package dto_request

type TokenBundleDTORequest struct {
	Name         string `json:"name" binding:"required"`
	TokenAmount  int    `json:"token_amount" binding:"required"`
	PriceRp      int    `json:"price_rp" binding:"required"`
	Label        string `json:"label,omitempty"`
	DisplayOrder int    `json:"display_order" binding:"required"`
	IsActive     *bool  `json:"is_active" binding:"required"`
}

type NewCustomPricingDTORequest struct {
	IsEnable                 *bool                         `json:"is_enable" binding:"required"`
	MinToken                 int                           `json:"min_token" binding:"required,gte=1,ltfield=MaxToken"`
	MaxToken                 int                           `json:"max_token" binding:"required,gte=1"`
	RecommendedPricePerToken int                           `json:"recommended_price_per_token" binding:"required,gte=1"`
	Tiers                    []CustomPricingTierDTORequest `json:"tiers" binding:"required,dive"`
}

type CustomPricingTierDTORequest struct {
	FromToken   int     `json:"from_token" binding:"required,gte=1,ltefield=ToToken"`
	ToToken     int     `json:"to_token" binding:"required,gte=1"`
	DiscountPct float64 `json:"discount_pct" binding:"required,gte=0,lte=100"`
}

type KPIDTORequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type CreateTokenTransactionDTORequest struct {
	TopupType     string `json:"topup_type" binding:"required,oneof=BUNDLE CUSTOM"`
	BundleID      string `json:"bundle_id,omitempty"`
	Amount        int    `json:"amount,omitempty" binding:"gte=1"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}
