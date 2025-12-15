package midtrans

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/payment_handler"

	"github.com/google/uuid"
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

func (m *MidtransClient) CreateMembershipPaymentRequest(grossAmount float64, planName string, email string) (string, string, error) {
	item := midtrans.ItemDetails{
		Name:     planName,
		Qty:      1,
		Price:    int64(grossAmount),
		Category: "membership",
	}

	customer := midtrans.CustomerDetails{
		Email: email,
	}

	transaction := snap.Request{
		Items:          &[]midtrans.ItemDetails{item},
		CustomerDetail: &customer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  uuid.NewString(),
			GrossAmt: int64(grossAmount),
		},
	}

	url, err := m.midtrans.CreateTransactionUrl(&transaction)
	if err != nil {
		return "", "", nil
	}

	return url, transaction.TransactionDetails.OrderID, nil
}

func (m *MidtransClient) CreateTokenPaymentRequest(req dto_request.MakeNewTransactionTokenRequest, email string, grossAmount float64) (string, string, error) {
	item := midtrans.ItemDetails{
		Name:     "Token",
		Price:    int64(grossAmount),
		Category: "token",
	}

	customer := midtrans.CustomerDetails{
		Email: email,
	}

	transaction := snap.Request{
		Items:          &[]midtrans.ItemDetails{item},
		CustomerDetail: &customer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  uuid.NewString(),
			GrossAmt: int64(grossAmount),
		},
	}

	url, err := m.midtrans.CreateTransactionUrl(&transaction)
	if err != nil {
		return "", "", nil
	}

	return url, transaction.TransactionDetails.OrderID, nil
}
