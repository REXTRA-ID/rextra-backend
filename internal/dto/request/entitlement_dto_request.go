package dto_request

type CreateEntitlementRequest struct {
	FeatureID        string  `json:"feature_id" binding:"required"`
	SubFeatureID     *string `json:"sub_feature_id,omitempty"`
	ActionCategoryID string  `json:"action_category_id" binding:"required"`
	Name             string  `json:"name" binding:"required,min=3,max=200"`
	Description      string  `json:"description"`
	Level            string  `json:"level" binding:"required,oneof=fitur sub_fitur"`
	RestrictionType  string  `json:"restriction_type" binding:"required,oneof=unlimited token_gated frequency_limited locked"`
	TokenCost        int     `json:"token_cost" binding:"min=0"`
}

type UpdateEntitlementRestrictionRequest struct {
	RestrictionType string `json:"restriction_type" binding:"required,oneof=unlimited token_gated frequency_limited locked"`
	TokenCost       int    `json:"token_cost" binding:"min=0"`
}
