package dto_request

type CreateFeatureRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Prefix      string `json:"prefix" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=tunggal bertingkat"`
}

type UpdateFeatureRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Prefix      string `json:"prefix" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=active inactive"`
	// Type tidak bisa diubah setelah dibuat — jika Feature sudah punya SubFeature
	// atau Entitlement, type tidak relevan untuk diubah dan rawan merusak data.
}
