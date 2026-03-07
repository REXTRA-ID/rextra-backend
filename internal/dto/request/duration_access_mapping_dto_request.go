package dto_request

type CreateDurationAccessMappingRequest struct {
	EntitlementID   string  `json:"entitlement_id" binding:"required"`
	RestrictionType string  `json:"restriction_type" binding:"required,oneof=unlimited token_gated frequency_limited locked"`
	TokenCost       int     `json:"token_cost"`
	UsageLimit      int     `json:"usage_limit"`
	ResetPeriod     *string `json:"reset_period,omitempty"`
}

type UpdateDurationAccessMappingRequest struct {
	RestrictionType string  `json:"restriction_type" binding:"required,oneof=unlimited token_gated frequency_limited locked"`
	TokenCost       int     `json:"token_cost"`
	UsageLimit      int     `json:"usage_limit"`
	ResetPeriod     *string `json:"reset_period,omitempty"`
	Status          string  `json:"status" binding:"required,oneof=aktif nonaktif"`
}
