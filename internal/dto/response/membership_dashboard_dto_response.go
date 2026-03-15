package dto_response

type MembershipDashboardResponse struct {
	Membership MembershipStatusInfo `json:"membership"`
	Wallet     TokenWalletInfo      `json:"wallet"`
	UIState    MembershipUIState    `json:"ui_state"`
	Benefits   []FeatureBenefitGroup `json:"benefits"`
}

type MembershipStatusInfo struct {
	PlanName       string  `json:"plan_name"`
	TierLabel      string  `json:"tier_label"`
	EmblemKey      string  `json:"emblem_key"`
	StatusLabel    string  `json:"status_label"` // e.g., "REXTRA CLUB"
	ActiveUntil    *string `json:"active_until,omitempty"`
	RemainingDays  int     `json:"remaining_days"`
}

type TokenWalletInfo struct {
	TokenBalance int64 `json:"token_balance"`
}

type MembershipUIState struct {
	IsNearExpiry      bool `json:"is_near_expiry"`
	IsExpired         bool `json:"is_expired"`
	HasClaimedStarter bool `json:"has_claimed_starter"`
	ThemeColor        string `json:"theme_color"` // Optional: HEX or Label
}

type FeatureBenefitGroup struct {
	FeatureName    string               `json:"feature_name"`
	Description    string               `json:"description"`
	IconKey        string               `json:"icon_key"`
	AccessType     string               `json:"access_type"` // Unlimited, Limited, Token, Varied
	AccessLabel    string               `json:"access_label"` // e.g., "Tanpa Batas"
	SubFeatures    []SubFeatureBenefit  `json:"sub_features,omitempty"`
}

type SubFeatureBenefit struct {
	Name        string `json:"name"`
	AccessType  string `json:"access_type"`
	AccessLabel string `json:"access_label"`
	QuotaLeft   *int   `json:"quota_left,omitempty"`
}
