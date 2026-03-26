package dto_request

type CheckoutPrepareRequest struct {
	PlanID     string `form:"plan_id" binding:"required"`
	ChangeType string `form:"change_type" binding:"required,oneof=PEMBELIAN_BARU RENEWAL UPGRADE DOWNGRADE"`
}

type CheckoutCalculateRequest struct {
	PlanID               string  `json:"plan_id" binding:"required"`
	DurationID           string  `json:"duration_id" binding:"required"`
	ChangeType           string  `json:"change_type" binding:"required,oneof=PEMBELIAN_BARU RENEWAL UPGRADE DOWNGRADE"`
	TokenBundlePackageID *string `json:"token_bundle_package_id"`
	PromoCode            *string `json:"promo_code"`
	UseCredit            bool    `json:"use_credit"`
}

type CheckoutInitiateRequest struct {
	PlanID               string  `json:"plan_id" binding:"required"`
	DurationID           string  `json:"duration_id" binding:"required"`
	ChangeType           string  `json:"change_type" binding:"required,oneof=PEMBELIAN_BARU RENEWAL UPGRADE DOWNGRADE"`
	TokenBundlePackageID *string `json:"token_bundle_package_id"`
	PromoCode            *string `json:"promo_code"`
	UseCredit            bool    `json:"use_credit"`
	PaymentMethod        string  `json:"payment_method" binding:"required"`
}

type CheckoutCancelRequest struct {
	CancelReason string  `json:"cancel_reason" binding:"required,oneof=CHANGE_ORDER BUDGET NOT_NOW PAYMENT_ISSUE OTHER"`
	CancelNote   *string `json:"cancel_note"`
}

type MyTransactionFilterRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=50"`
	Status   string `form:"status" binding:"omitempty,oneof=pending paid failed expired cancelled"`
}
