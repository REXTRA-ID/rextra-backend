package dto_request

type MakeNewTransactionRequest struct {
	PaymentType   string  `json:"payment_type"`
	PlanID        *string `json:"plan_id,omitempty"`
	DurationID    *string `json:"duration_id,omitempty"`
	TokenQuantity *int    `json:"token_quantity,omitempty"`
	GrossAmount   float64 `json:"gross_amount"`
}
