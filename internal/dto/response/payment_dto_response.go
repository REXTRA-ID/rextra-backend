package dto_response

type MakeNewTransactionResponse struct {
	RedirectURL string `json:"payment_url"`
	Message     string `json:"message"`
}

type GetPaymentTransactionsResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	PaymentType   string  `json:"payment_type"`
	GrossAmount   float64 `json:"gross_amount"`
	FinalAmount   float64 `json:"final_amount"`
	PromoCode     string  `json:"promo_code,omitempty"`
	PaymentMethod string  `json:"paymont_method"`
	PaymentStatus string  `json:"payment_status"`
	PaidAt        string  `json:"paid_at,omitempty"`
	ExpiredAt     string  `json:"expired_at,omitempty"`
}
