package dto_request

type MakeNewTransactionTokenRequest struct {
	PaymentType   string  `json:"payment_type"`
	TokenQuantity int     `json:"token_quantity"`
	GrossAmount   float64 `json:"gross_amount"`
	Description   string  `json:"description"`
	PromoCode     *string `json:"promo_code"`
}

type MakeNewTransactionMembershipRequest struct {
	PaymentType string  `json:"payment_type"`
	PlanID      string  `json:"plan_id"`
	DurationID  string  `json:"duration_id"`
	Description string  `json:"description"`
	PromoCode   *string `json:"promo_code"`
}

type PaymentWebhookRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
}

type MidtransCallback struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	SettlementTime    string `json:"settlement_time"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}
