package dto_request

type MakeNewTransactionRequest struct {
	PaymentType   string  `json:"payment_type"`
	PlanID        *string `json:"plan_id,omitempty"`
	DurationID    *string `json:"duration_id,omitempty"`
	TokenQuantity *int    `json:"token_quantity,omitempty"`
	GrossAmount   float64 `json:"gross_amount"`
}

type MakeNewTransactionTokenRequest struct {
	PaymentType   string  `json:"payment_type"`
	TokenQuantity int     `json:"token_quantity"`
	GrossAmount   float64 `json:"gross_amount"`
	Description   string  `json:"description"`
}

type MakeNewTransactionMembershipRequest struct {
	PaymentType string  `json:"payment_type"`
	PlanID      string  `json:"plan_id"`
	DurationID  string  `json:"duration_id"`
	GrossAmount float64 `json:"gross_amount"`
	Description string  `json:"description"`
}
