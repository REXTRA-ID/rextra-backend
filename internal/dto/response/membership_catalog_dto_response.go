package dto_response

type MembershipCatalogResponse struct {
	CurrentContext *CurrentSubscriptionContext `json:"current_context,omitempty"`
	Plans          []PlanCatalogItem           `json:"plans"`
}

type CurrentSubscriptionContext struct {
	Message        string `json:"message"` // e.g., "Kamu sedang berlangganan Pro Plan — 1 bulan."
	CurrentPlanID  string `json:"current_plan_id"`
	CurrentPlanName string `json:"current_plan_name"`
}

type PlanCatalogItem struct {
	ID             string               `json:"id"`
	PlanName       string               `json:"plan_name"`
	TierLabel      string               `json:"tier_label"`
	EmblemKey      string               `json:"emblem_key"`
	ThemeColor     string               `json:"theme_color"`
	MarketingIntro string               `json:"marketing_intro"`
	PricingInfo    string               `json:"pricing_info"` // e.g., "Mulai dari Rp10.000"
	TokenBonusInfo string               `json:"token_bonus_info"` // e.g., "100 Token"
	Label          string               `json:"label,omitempty"` // e.g., "Paling Direkomendasikan"
	
	CTA            PlanCTAInfo          `json:"cta"`
	Benefits       []FeatureBenefitGroup `json:"benefits"`
}

type PlanCTAInfo struct {
	Label      string `json:"label"`       // e.g., "Upgrade ke plan ini"
	ChangeType string `json:"change_type"` // RENEWAL, UPGRADE, DOWNGRADE, PEMBELIAN_BARU
	IsCurrent  bool   `json:"is_current"`
}
