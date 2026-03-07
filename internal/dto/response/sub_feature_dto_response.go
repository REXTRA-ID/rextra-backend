package dto_response

type GetSubFeatureResponse struct {
	Id               string `json:"id"`
	FeatureId        string `json:"feature_id"`
	FeatureName      string `json:"feature_name"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	EntitlementCount int    `json:"entitlement_count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
