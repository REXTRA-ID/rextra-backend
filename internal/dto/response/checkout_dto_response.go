package dto_response

type CheckoutPrepareResponse struct {
	Plan                 CheckoutPlanInfo          `json:"plan"`
	CurrentMembership    *CheckoutMembershipInfo   `json:"current_membership,omitempty"`
	EligibleVoucherCount int                       `json:"eligible_voucher_count"`
	TokenBundles         []CheckoutTokenBundleInfo `json:"token_bundles"`
}

type CheckoutPlanInfo struct {
	ID           string                 `json:"id"`
	PlanName     string                 `json:"plan_name"`
	TierLabel    string                 `json:"tier_label"`
	EmblemKey    string                 `json:"emblem_key"`
	DurationMode string                 `json:"duration_mode"`
	Durations    []CheckoutDurationInfo `json:"durations"`
}

type CheckoutDurationInfo struct {
	ID             string  `json:"id"`
	DurationMonths int     `json:"duration_months"`
	Price          int64   `json:"price"`
	FinalPrice     int64   `json:"final_price"`
	DiscountPct    float64 `json:"discount_pct"`
	PricePerMonth  int64   `json:"price_per_month"`
}

type CheckoutMembershipInfo struct {
	PlanName        string `json:"plan_name"`
	DurationMonths  int    `json:"duration_months"`
	RemainingDays   int    `json:"remaining_days"`
	EstimatedCredit int64  `json:"estimated_credit"`
}

type CheckoutTokenBundleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TokenAmount int64  `json:"token_amount"`
	PriceRp     int64  `json:"price_rp"`
	Label       string `json:"label,omitempty"`
}

type CheckoutCalculateResponse struct {
	MembershipPrice int64  `json:"membership_price"`
	TokenPrice      int64  `json:"token_price"`
	SubtotalAmount  int64  `json:"subtotal_amount"`
	DiscountAmount  int64  `json:"discount_amount"`
	DurationCredit  int64  `json:"duration_credit"`
	TotalAmount     int64  `json:"total_amount"`
	PromoApplied    bool   `json:"promo_applied"`
	PromoCode       string `json:"promo_code,omitempty"`
	PromoErrMsg     string `json:"promo_err_msg,omitempty"`
	RemainingDays   int    `json:"remaining_days"`
	CreditAvailable int64  `json:"credit_available"`
}

type CheckoutInitiateResponse struct {
	TransactionID string `json:"transaction_id"`
	PaymentURL    string `json:"payment_url"`
	PayCode       string `json:"pay_code,omitempty"`
	ExpiresAt     string `json:"expires_at"`
	TotalAmount   int64  `json:"total_amount"`
}

type CheckoutRepeatResponse struct {
	PlanID               *string `json:"plan_id,omitempty"`
	DurationID           *string `json:"duration_id,omitempty"`
	ChangeType           string  `json:"change_type"`
	TokenBundlePackageID *string `json:"token_bundle_package_id,omitempty"`
	PromoCode            *string `json:"promo_code,omitempty"`
	PromoExpired         bool    `json:"promo_expired"`
}

type MyTransactionListResponse struct {
	ID            string                     `json:"id"`
	TransactionID string                     `json:"transaction_id"`
	ChangeType    string                     `json:"change_type"`
	PrimaryStatus string                     `json:"primary_status"`
	PaymentStatus string                     `json:"payment_status"`
	TotalAmount   int64                      `json:"total_amount"`
	Items         []MyTransactionItemPreview `json:"items"`
	CreatedAt     string                     `json:"created_at"`
	PaidAt        string                     `json:"paid_at,omitempty"`
}

type MyTransactionItemPreview struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MyTransactionDetailResponse struct {
	ID             string                    `json:"id"`
	TransactionID  string                    `json:"transaction_id"`
	ChangeType     string                    `json:"change_type"`
	FromPlan       string                    `json:"from_plan,omitempty"`
	ToPlan         string                    `json:"to_plan"`
	FromDuration   int                       `json:"from_duration_months,omitempty"`
	ToDuration     int                       `json:"to_duration_months"`
	PrimaryStatus  string                    `json:"primary_status"`
	PaymentStatus  string                    `json:"payment_status"`
	PaymentMethod  string                    `json:"payment_method,omitempty"`
	PaymentURL     string                    `json:"payment_url,omitempty"`
	PayCode        string                    `json:"pay_code,omitempty"`
	SubtotalAmount int64                     `json:"subtotal_amount"`
	DiscountAmount int64                     `json:"discount_amount"`
	DurationCredit int64                     `json:"duration_credit"`
	AdminFee       int64                     `json:"admin_fee"`
	TotalAmount    int64                     `json:"total_amount"`
	PromoCode      string                    `json:"promo_code,omitempty"`
	Items          []MyTransactionItemDetail `json:"items"`
	CancelReason   string                    `json:"cancel_reason,omitempty"`
	CancelNote     string                    `json:"cancel_note,omitempty"`
	CreatedAt      string                    `json:"created_at"`
	PaidAt         string                    `json:"paid_at,omitempty"`
	ExpiresAt      string                    `json:"expires_at,omitempty"`
	CanceledAt     string                    `json:"canceled_at,omitempty"`
	IsCancellable  bool                      `json:"is_cancellable"`
	IsRepeatable   bool                      `json:"is_repeatable"`
}

type MyTransactionItemDetail struct {
	Type            string `json:"type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Price           int64  `json:"price"`
	Quantity        int    `json:"quantity"`
	Subtotal        int64  `json:"subtotal"`
	TokenAmount     int    `json:"token_amount,omitempty"`
	BundlePackageID string `json:"bundle_package_id,omitempty"`
}

type PaymentChannelResponse struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	IconURL  string `json:"icon_url"`
	AdminFee int64  `json:"admin_fee"`
}
