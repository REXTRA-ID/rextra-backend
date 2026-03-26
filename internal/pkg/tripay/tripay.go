package tripay

import (
	"os"

	e "rextra-backend/internal/entity"
	u "rextra-backend/internal/utils"

	"github.com/zakirkun/go-tripay/client"
	"github.com/zakirkun/go-tripay/utils"
)

type TripayOrderItem struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreatePaymentRequest struct {
	Method        string
	MerchantRef   string
	Amount        int64
	CustomerName  string
	CustomerEmail string
	OrderItems    []TripayOrderItem
	ExpiredHours  int
}

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

func (c *TripayClient) CreatePaymentTransaction(req CreatePaymentRequest) (e.TransactionResponse, error) {
	signStr := utils.Signature{
		Amount:       req.Amount,
		PrivateKey:   c.client.PrivateKey,
		MerchantCode: c.client.MerchantCode,
		MerchanReff:  req.MerchantRef,
	}

	c.client.SetSignature(signStr)

	var tripayOrderItems []client.OrderItemClosePaymentRequest
	for _, item := range req.OrderItems {
		tripayOrderItems = append(tripayOrderItems, client.OrderItemClosePaymentRequest{
			SKU:        item.SKU,
			Name:       item.Name,
			Price:      item.Price,
			Quantity:   item.Quantity,
			ProductURL: "",
			ImageURL:   "",
		})
	}

	bodyReq := client.ClosePaymentBodyRequest{
		Method:        utils.TRIPAY_CHANNEL(req.Method),
		MerchantRef:   req.MerchantRef,
		Amount:        int(req.Amount),
		CustomerEmail: req.CustomerEmail,
		CustomerName:  req.CustomerName,
		ExpiredTime:   client.SetTripayExpiredTime(req.ExpiredHours),
		Signature:     c.client.GetSignature(),
		OrderItems:    tripayOrderItems,
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
		PayCode:      response.Data.PayCode,
	}, nil
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
