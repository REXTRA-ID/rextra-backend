package payment_handler

import (
	dto_request "rextra-backend/internal/dto/request"
)

type PaymentService interface {
	CreatePaymentRequest(req dto_request.MakeNewTransactionRequest) (any, error)
}
