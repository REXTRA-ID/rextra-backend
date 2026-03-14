package tripay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TripayClient struct {
	MerchantCode string
	ApiKey       string
	PrivateKey   string
	BaseUrl      string
}

type TripayOrderItem struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreatePaymentRequest struct {
	Method         string            `json:"method"`
	MerchantRef    string            `json:"merchant_ref"`
	Amount         int64             `json:"amount"`
	CustomerName   string            `json:"customer_name"`
	CustomerEmail  string            `json:"customer_email"`
	CustomerPhone  string            `json:"customer_phone"`
	OrderItems     []TripayOrderItem `json:"order_items"`
	CallbackUrl    string            `json:"callback_url"`
	ReturnUrl      string            `json:"return_url"`
	ExpiredTime    int64             `json:"expired_time"`
	Signature      string            `json:"signature"`
}

type TripayResponse struct {
	Reference      string `json:"reference"`
	MerchantRef    string `json:"merchant_ref"`
	PaymentMethod  string `json:"payment_method"`
	PaymentName    string `json:"payment_name"`
	CustomerName   string `json:"customer_name"`
	CustomerEmail  string `json:"customer_email"`
	Amount         int64  `json:"amount"`
	FeeMerchant    int64  `json:"fee_merchant"`
	FeeCustomer    int64  `json:"fee_customer"`
	TotalFee       int64  `json:"total_fee"`
	AmountReceived int64  `json:"amount_received"`
	CheckoutURL    string `json:"checkout_url"`
	Status         string `json:"status"`
	PayCode        string `json:"pay_code"`
	ExpiredTime    int64  `json:"expired_time"`
}

type PaymentChannel struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Group   string `json:"group"`
	TotalFee struct {
		Flat    int64  `json:"flat"`
		Percent string `json:"percent"`
	} `json:"total_fee"`
	Active  bool   `json:"active"`
}

func NewTripayClient() TripayClient {
	baseUrl := "https://tripay.co.id/api-sandbox"
	if os.Getenv("TRIPAY_MODE") == "live" {
		baseUrl = "https://tripay.co.id/api"
	}

	return TripayClient{
		MerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		ApiKey:       os.Getenv("TRIPAY_API_KEY"),
		PrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		BaseUrl:      baseUrl,
	}
}

func (t *TripayClient) CreatePaymentTransaction(req CreatePaymentRequest) (TripayResponse, error) {
	signaturePayload := fmt.Sprintf("%s%s%d", t.MerchantCode, req.MerchantRef, req.Amount)
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write([]byte(signaturePayload))
	req.Signature = hex.EncodeToString(h.Sum(nil))

	jsonReq, _ := json.Marshal(req)
	client := &http.Client{Timeout: 15 * time.Second}
	
	url := fmt.Sprintf("%s/transaction/create", t.BaseUrl)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+t.ApiKey)

	resp, err := client.Do(httpReq)
	if err != nil { return TripayResponse{}, err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var fullResp struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &fullResp); err != nil {
		return TripayResponse{}, fmt.Errorf("failed to parse tripay response: %s", string(body))
	}
	if !fullResp.Success {
		return TripayResponse{}, fmt.Errorf("tripay error: %s", fullResp.Message)
	}
	var data TripayResponse
	json.Unmarshal(fullResp.Data, &data)
	return data, nil
}

func (t *TripayClient) VerifyCallbackSignature(payload []byte, remoteSignature string) bool {
	h := hmac.New(sha256.New, []byte(t.PrivateKey))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return expectedSignature == remoteSignature
}

func (t *TripayClient) GetPaymentChannels() ([]PaymentChannel, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/merchant/payment-channel", t.BaseUrl)
	
	httpReq, _ := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("Authorization", "Bearer "+t.ApiKey)

	resp, err := client.Do(httpReq)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[DEBUG] Tripay GetPaymentChannels Raw: %s", string(body))

	var rawResult struct {
		Success bool `json:"success"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(body, &rawResult)

	var res []PaymentChannel
	for _, item := range rawResult.Data {
		ch := PaymentChannel{
			Code: fmt.Sprintf("%v", item["code"]),
			Name: fmt.Sprintf("%v", item["name"]),
			Group: fmt.Sprintf("%v", item["group"]),
			Active: true,
		}
		if tf, ok := item["total_fee"].(map[string]interface{}); ok {
			ch.TotalFee.Flat = toInt64(tf["flat"])
			ch.TotalFee.Percent = fmt.Sprintf("%v", tf["percent"])
		}
		res = append(res, ch)
	}
	return res, nil
}

func (t *TripayClient) GetFeeCalculation(amount int64, code string) ([]PaymentChannel, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/merchant/fee-calculator?amount=%d&code=%s", t.BaseUrl, amount, code)
	
	httpReq, _ := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("Authorization", "Bearer "+t.ApiKey)

	resp, err := client.Do(httpReq)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("[DEBUG] Tripay GetFeeCalculation Raw: %s", string(body))

	var rawResult struct {
		Success bool `json:"success"`
		Data []map[string]interface{} `json:"data"`
	}
	json.Unmarshal(body, &rawResult)

	var res []PaymentChannel
	for _, item := range rawResult.Data {
		ch := PaymentChannel{
			Code: fmt.Sprintf("%v", item["code"]),
			Name: fmt.Sprintf("%v", item["name"]),
		}
		if fee, ok := item["fee"].(map[string]interface{}); ok {
			ch.TotalFee.Flat = toInt64(fee["flat"])
			ch.TotalFee.Percent = fmt.Sprintf("%v", fee["percent"])
		}
		res = append(res, ch)
	}
	return res, nil
}

func toInt64(v interface{}) int64 {
	if v == nil { return 0 }
	switch val := v.(type) {
	case float64: return int64(val)
	case int64: return val
	case int: return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	}
	return 0
}
