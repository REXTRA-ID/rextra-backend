package tripay

import (
	"os"

	e "rextra-backend/internal/entity"
	u "rextra-backend/internal/utils"

	"github.com/zakirkun/go-tripay/client"
	"github.com/zakirkun/go-tripay/utils"
)

type TripayClient struct {
	client client.Client
}

func NewTripayClient() TripayClient {
	apiKey := os.Getenv("API_KEY")
	privateKey := os.Getenv("PRIVATE_KEY")
	merchantCode := os.Getenv("MERCHANT_CODE")

	return TripayClient{
		client: client.Client{
			MerchantCode: merchantCode,
			ApiKey:       apiKey,
			PrivateKey:   privateKey,
			Mode:         utils.MODE_DEVELOPMENT,
		},
	}
}

func (c *TripayClient) GetPaymentChannel() ([]e.PaymentChannel, error) {
	response, err := c.client.MerchantPay()
	if err != nil {
		return nil, err
	}
	var channels []e.PaymentChannel
	for _, channel := range response.Data {
		addChannel := e.PaymentChannel{
			Code:    channel.Code,
			Name:    channel.Name,
			IconUrl: channel.IconURL,
		}
		channels = append(channels, addChannel)
	}
	return channels, nil
}

func (c *TripayClient) CreatePayment(amount int64, method, customerEmail string, paymentType e.PaymentType) (e.TransactionResponse, error) {

	merchantReff := u.GenerateMerchantReff()

	signStr := utils.Signature{
		Amount:       amount,
		PrivateKey:   c.client.PrivateKey,
		MerchantCode: c.client.MerchantCode,
		MerchanReff:  merchantReff,
	}

	c.client.SetSignature(signStr)

	var sku, name string

	if paymentType == e.MEMBERSHIP {
		sku = "membership"
		name = "Beli membership"
	} else {
		sku = "token"
		name = "Beli token"
	}

	bodyReq := client.ClosePaymentBodyRequest{
		Method:        utils.TRIPAY_CHANNEL(method),
		MerchantRef:   merchantReff,
		Amount:        int(amount),
		CustomerEmail: customerEmail,
		ExpiredTime:   client.SetTripayExpiredTime(2),
		Signature:     c.client.GetSignature(),
		OrderItems: []client.OrderItemClosePaymentRequest{
			{
				SKU:        sku,
				Name:       name,
				Price:      int(amount),
				ProductURL: "",
				ImageURL:   "",
			},
		},
	}

	response, err := c.client.ClosePaymentRequestTransaction(bodyReq)
	if err != nil {
		return e.TransactionResponse{}, err
	}

	return e.TransactionResponse{
		Reference:    response.Data.Reference,
		MerchantReff: response.Data.MerchantRef,
		Status:       response.Data.Status,
		CheckoutUrl:  response.Data.CheckoutURL,
	}, nil
}
