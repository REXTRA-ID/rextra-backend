package dto_response

type GetEntitlementResponse struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Status      string `json:"status"`

	MappingCount int `json:"mapping_count"`

	FeatureId          string  `json:"feature_id"`
	FeatureName        string  `json:"feature_name"`
	FeaturePrefix      string  `json:"feature_prefix"`
	SubFeatureId       *string `json:"sub_feature_id,omitempty"`
	SubFeatureName     *string `json:"sub_feature_name,omitempty"`
	ActionCategoryId   string  `json:"action_category_id"`
	ActionCategoryName string  `json:"action_category_name"`
	ActionCategorySlug string  `json:"action_category_slug"`

	RestrictionType string `json:"restriction_type"`
	TokenCost       int    `json:"token_cost"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
