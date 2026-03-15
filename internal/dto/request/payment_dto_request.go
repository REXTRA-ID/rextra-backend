package dto_request

type MakeNewTransactionMembershipRequest struct {
	PlanID        string  `json:"plan_id" binding:"required"`
	DurationID    string  `json:"duration_id" binding:"required"`
	ChangeType    string  `json:"change_type" binding:"required,oneof=PEMBELIAN_BARU RENEWAL UPGRADE DOWNGRADE"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	PromoCode     *string `json:"promo_code"`
}

type CalculatePriceRequest struct {
	PlanID     string  `json:"plan_id" binding:"required"`
	DurationID string  `json:"duration_id" binding:"required"`
	ChangeType string  `json:"change_type" binding:"required,oneof=PEMBELIAN_BARU RENEWAL UPGRADE DOWNGRADE"`
	PromoCode  *string `json:"promo_code"`
}

type PaymentTransactionFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   string `form:"status" binding:"omitempty,oneof=pending paid failed expired cancelled"`
	PlanName string `form:"plan_name" binding:"omitempty"`
	UserID   string `form:"user_id" binding:"omitempty"`
	DateFrom string `form:"date_from" binding:"omitempty"`
	DateTo   string `form:"date_to" binding:"omitempty"`
	Search   string `form:"search" binding:"omitempty"`
}

type MyTransactionFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   string `form:"status" binding:"omitempty,oneof=pending paid failed expired cancelled"`
}

type CancelTransactionRequest struct {
	Note string `json:"note" binding:"required,min=5"`
}

type TripayCallbackRequest struct {
	Reference     string  `json:"reference"`
	MerchantRef   string  `json:"merchant_ref"`
	PaymentMethod string  `json:"payment_method"`
	TotalAmount   int64   `json:"total_amount"`
	Status        string  `json:"status"`
	PaidAt        *int64  `json:"paid_at"`
	Note          *string `json:"note"`
}
