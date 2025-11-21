package midtrans

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/payment_handler"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	midtrans *snap.Client
}

func NewMidtransClient() payment_handler.PaymentService {
	var client snap.Client
	client.New("MIDTRANS_SERVER_KEY", midtrans.Sandbox)
	return &MidtransClient{midtrans: &client}
}

func (m *MidtransClient) CreatePaymentRequest(req dto_request.MakeNewTransactionRequest) (any, error) {
	// implement nanti
	return 0, nil
}
