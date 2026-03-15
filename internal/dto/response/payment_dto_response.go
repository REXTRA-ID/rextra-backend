package dto_response

type MakeNewTransactionResponse struct {
	TransactionID  string  `json:"transaction_id"`
	PaymentURL     *string `json:"payment_url,omitempty"`
	PayCode        *string `json:"pay_code,omitempty"`
	TotalAmount    int64   `json:"total_amount"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
	Message        string  `json:"message,omitempty"`
	SubtotalAmount int64   `json:"subtotal_amount"`
	DiscountAmount int64   `json:"discount_amount"`
	DurationCredit int64   `json:"duration_credit"`
	PromoCode      *string `json:"promo_code,omitempty"`
}

type GetPaymentTransactionListResponse struct {
	ID               string  `json:"id"`
	TransactionID    string  `json:"transaction_id"`
	UserID           string  `json:"user_id"`
	UserName         string  `json:"user_name"`
	UserEmail        string  `json:"user_email"`
	ToPlan           string  `json:"to_plan"`
	ToDurationMonths int     `json:"to_duration_months"`
	ChangeType       string  `json:"change_type"`
	TotalAmount      int64   `json:"total_amount"`
	PaymentMethod    *string `json:"payment_method,omitempty"`
	PaymentStatus    string  `json:"payment_status"`
	PrimaryStatus    string  `json:"primary_status"`
	PaymentURL       *string `json:"payment_url,omitempty"`
	PayCode          *string `json:"pay_code,omitempty"`
	CreatedAt        string  `json:"created_at"`
	PaidAt           *string `json:"paid_at,omitempty"`
	ExpiresAt        *string `json:"expires_at,omitempty"`
}

type GetPaymentTransactionDetailResponse struct {
	ID                 string  `json:"id"`
	TransactionID      string  `json:"transaction_id"`
	UserID             string  `json:"user_id"`
	UserName           string  `json:"user_name"`
	UserEmail          string  `json:"user_email"`
	ChangeType         string  `json:"change_type"`
	FromPlan           *string `json:"from_plan,omitempty"`
	ToPlan             string  `json:"to_plan"`
	FromDurationMonths *int    `json:"from_duration_months,omitempty"`
	ToDurationMonths   int     `json:"to_duration_months"`
	SubtotalAmount     int64   `json:"subtotal_amount"`
	DiscountAmount     int64   `json:"discount_amount"`
	DurationCredit     int64   `json:"duration_credit"`
	TotalAmount        int64   `json:"total_amount"`
	PromoCode          *string `json:"promo_code,omitempty"`
	PaymentProvider    string  `json:"payment_provider"`
	PaymentMethod      *string `json:"payment_method,omitempty"`
	PaymentStatus      string  `json:"payment_status"`
	PrimaryStatus      string  `json:"primary_status"`
	PaymentURL         *string `json:"payment_url,omitempty"`
	PayCode            *string `json:"pay_code,omitempty"`
	PaymentExternalID  string  `json:"payment_external_id"`
	CancelReason       *string `json:"cancel_reason,omitempty"`
	CancelNote         *string `json:"cancel_note,omitempty"`
	CanceledAt         *string `json:"canceled_at,omitempty"`
	CreatedAt          string  `json:"created_at"`
	PaidAt             *string `json:"paid_at,omitempty"`
	ExpiresAt          *string `json:"expires_at,omitempty"`
}

type PaginatedPaymentTransactionResponse struct {
	Data       []GetPaymentTransactionListResponse `json:"data"`
	Total      int                                 `json:"total"`
	Page       int                                 `json:"page"`
	PageSize   int                                 `json:"page_size"`
	TotalPages int                                 `json:"total_pages"`
}

type GetPaymentChannelResponse struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Group         string `json:"group"`
	IconURL       string `json:"icon_url"`
	FeeFlat       int64  `json:"fee_flat"`
	FeePercent    string `json:"fee_percent"`
	MinimumAmount int64  `json:"minimum_amount"`
	MaximumAmount int64  `json:"maximum_amount"`
}

type PriceCalculationResponse struct {
	SubtotalAmount int64   `json:"subtotal_amount"`
	DiscountAmount int64   `json:"discount_amount"`
	PromoCode      *string `json:"promo_code,omitempty"`
	DurationCredit int64   `json:"duration_credit"`
	TotalAmount    int64   `json:"total_amount"`
	PlanName       string  `json:"plan_name"`
	DurationMonths int     `json:"duration_months"`
	TokenAmount    int     `json:"token_amount"`
	BonusToken     int     `json:"bonus_token"`
}
