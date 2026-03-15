package dto_request

type CreateDurationAccessMappingRequest struct {
	EntitlementID string `json:"entitlement_id" binding:"required"`
	UsageLimit int `json:"usage_limit" binding:"min=0"`
}

type UpdateDurationAccessMappingRequest struct {
	Status string `json:"status" binding:"required,oneof=aktif nonaktif"`
	UsageLimit *int `json:"usage_limit,omitempty" binding:"omitempty,min=0"`
}
