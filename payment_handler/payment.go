package payment_handler

import (
	dto_request "rextra-backend/internal/dto/request"
)

type PaymentService interface {
	CreateMembershipPaymentRequest(grossAmount float64, planName string, email string) (string, string, error)
	CreateTokenPaymentRequest(req dto_request.MakeNewTransactionTokenRequest, email string) (string, string, error)
}
