package tripay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type InstructionStep struct {
	Title string   `json:"title"`
	Steps []string `json:"steps"`
}

type OrderItem map[string]interface{}

type TransactionPayload struct {
	Method        string      `json:"method"`
	MerchantRef   string      `json:"merchant_ref"`
	Amount        int         `json:"amount"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	CustomerPhone string      `json:"customer_phone"`
	OrderItems    []OrderItem `json:"order_items"`
	ReturnURL     string      `json:"return_url"`
	ExpiredTime   int64       `json:"expired_time"`
	Signature     string      `json:"signature"`
}

type ResponseOrderItem struct {
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	Subtotal   int    `json:"subtotal"`
	ProductURL string `json:"product_url"`
	ImageURL   string `json:"image_url"`
}

type TransactionData struct {
	Reference            string              `json:"reference"`
	MerchantRef          string              `json:"merchant_ref"`
	PaymentSelectionType string              `json:"payment_selection_type"`
	PaymentMethod        string              `json:"payment_method"`
	PaymentName          string              `json:"payment_name"`
	CustomerName         string              `json:"customer_name"`
	CustomerEmail        string              `json:"customer_email"`
	CustomerPhone        string              `json:"customer_phone"`
	CallbackURL          string              `json:"callback_url"`
	ReturnURL            string              `json:"return_url"`
	Amount               int                 `json:"amount"`
	FeeMerchant          int                 `json:"fee_merchant"`
	FeeCustomer          int                 `json:"fee_customer"`
	TotalFee             int                 `json:"total_fee"`
	AmountReceived       int                 `json:"amount_received"`
	PayCode              string              `json:"pay_code"`
	PayURL               *string             `json:"pay_url"` // pakai pointer karena bisa null
	CheckoutURL          string              `json:"checkout_url"`
	Status               string              `json:"status"`
	ExpiredTime          int64               `json:"expired_time"`
	OrderItems           []ResponseOrderItem `json:"order_items"`
	Instructions         []InstructionStep   `json:"instructions"`
	QRString             *string             `json:"qr_string"` // pakai pointer karena bisa null
	QRURL                *string             `json:"qr_url"`    // pakai pointer karena bisa null
}

func GetInstruction(code string) ([]InstructionStep, error) {
	apiKey := os.Getenv("TRIPAY_API_KEY")
	url := "https://tripay.co.id/api-sandbox/payment/instruction?code=" + code

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Data    []InstructionStep `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return result.Data, nil
}

func GenerateSignature(privateKey, merchantCode, merchantRef string, amount int) string {
	data := merchantCode + merchantRef + fmt.Sprintf("%d", amount)
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateTransaction(payload TransactionPayload, orderItems []OrderItem) (TransactionData, error) {
	if err := validateOrderItems(orderItems); err != nil {
		return TransactionData{}, err
	}

	apiKey := os.Getenv("TRIPAY_API_KEY")
	privateKey := os.Getenv("TRIPAY_PRIVATE_KEY")
	merchantCode := os.Getenv("TRIPAY_MERCHANT_CODE")

	ExpiredTime := time.Now().Unix() + (24 * 60 * 60)
	payload.ExpiredTime = ExpiredTime
	payload.Signature = GenerateSignature(privateKey, merchantCode, payload.MerchantRef, payload.Amount)
	payload.OrderItems = orderItems
	payload.ReturnURL = os.Getenv("TRIPAY_RETURN_URL")

	body, err := json.Marshal(payload)
	if err != nil {
		return TransactionData{}, err

	}

	req, err := http.NewRequest("POST", "https://tripay.co.id/api-sandbox/transaction/create", bytes.NewBuffer(body))
	if err != nil {
		return TransactionData{}, err

	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TransactionData{}, err

	}
	defer resp.Body.Close()

	var result struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    TransactionData `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return TransactionData{}, err

	}

	if !result.Success {
		return TransactionData{}, errors.New(result.Message)
	}

	result.Data.ExpiredTime = ExpiredTime

	return result.Data, nil
}

func validateOrderItems(items []OrderItem) error {
	for i, item := range items {
		if _, ok := item["name"]; !ok {
			return fmt.Errorf("order_items[%d]: 'name' wajib diisi", i)
		}
		if _, ok := item["price"]; !ok {
			return fmt.Errorf("order_items[%d]: 'price' wajib diisi", i)
		}
		if _, ok := item["quantity"]; !ok {
			return fmt.Errorf("order_items[%d]: 'quantity' wajib diisi", i)
		}
	}
	return nil
}
