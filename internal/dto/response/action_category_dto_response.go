package dto_response

type GetActionCategoryResponse struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	EntitlementCount int    `json:"entitlement_count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
