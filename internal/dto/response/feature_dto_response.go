package dto_response

type GetFeatureResponse struct {
	Id               string                  `json:"id"`
	Name             string                  `json:"name"`
	Slug             string                  `json:"slug"`
	Prefix           string                  `json:"prefix"`
	Description      string                  `json:"description"`
	Type             string                  `json:"type"`
	Status           string                  `json:"status"`
	EntitlementCount int                     `json:"entitlement_count"`
	SubFeatures      []GetSubFeatureResponse `json:"sub_features,omitempty"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
}
