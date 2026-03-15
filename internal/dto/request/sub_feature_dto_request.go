package dto_request

type CreateSubFeatureRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

type UpdateSubFeatureRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=active inactive"`
}
