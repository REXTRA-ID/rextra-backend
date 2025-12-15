package xendit

import (
	"os"
	dto_request "rextra-backend/internal/dto/request"
	payment_handler "rextra-backend/payment_handler"

	"github.com/xendit/xendit-go/v7"
)

type XenditService struct {
	xnd *xendit.APIClient
}

func NewXenditService() payment_handler.PaymentService {
	xnd := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))
	return &XenditService{xnd: xnd}
}

func (c *XenditService) CreateMembershipPaymentRequest(grossAmount float64, planName string, email string) (string, string, error) {
	return "", "", nil
}

func (c *XenditService) CreateTokenPaymentRequest(req dto_request.MakeNewTransactionTokenRequest, email string, grossAmount float64) (string, string, error) {
	return "", "", nil
}
