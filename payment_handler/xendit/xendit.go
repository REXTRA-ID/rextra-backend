package xendit

import (
	"context"
	"os"
	dto_request "rextra-backend/internal/dto/request"
	payment_handler "rextra-backend/payment_handler"

	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7"
)

type XenditService struct {
	xnd *xendit.APIClient
}

func NewXenditService() payment_handler.PaymentService {
	xnd := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))
	return &XenditService{xnd: xnd}
}

func (c *XenditService) CreatePaymentRequest(req dto_request.MakeNewTransactionRequest) (any, error) {
	request := c.xnd.PaymentRequestApi.GetPaymentRequestByID(context.Background(), uuid.New().String())
	return request, nil
}
