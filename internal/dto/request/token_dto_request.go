package dto_request

type TokenBundleDTORequest struct {
	Name         string `json:"name" binding:"required"`
	TokenAmount  int    `json:"token_amount" binding:"required"`
	PriceRp      int    `json:"price_rp" binding:"required"`
	Label        string `json:"label,omitempty"`
	DisplayOrder int    `json:"display_order" binding:"required"`
	IsActive     *bool  `json:"is_active" binding:"required"`
}
