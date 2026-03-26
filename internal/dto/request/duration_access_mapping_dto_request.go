package dto_request

type CreateDurationAccessMappingRequest struct {
	EntitlementID string `json:"entitlement_id" binding:"required"`
	// UsageLimit hanya relevan jika entitlement.RestrictionType = "frequency_limited"
	// Nilainya ditentukan manual per plan per durasi oleh admin
	UsageLimit int `json:"usage_limit"`
}

type UpdateDurationAccessMappingRequest struct {
	// UsageLimit hanya relevan jika entitlement.RestrictionType = "frequency_limited"
	UsageLimit int    `json:"usage_limit"`
	Status     string `json:"status" binding:"required,oneof=aktif nonaktif"`
}
