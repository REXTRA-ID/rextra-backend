package dto_request

type CreateFeatureRequest struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Status   string  `json:"status"`
	ParentId *string `json:"parent_id,omitempty"`
}
