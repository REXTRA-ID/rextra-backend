package dto_response

type GetFeatureResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
	Type   string `json:"type"`
}
