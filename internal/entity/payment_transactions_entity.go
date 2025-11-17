package entity

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type PaymentStatus string

const (
	VAPAYMENT      PaymentStatus = "va"
	EWALLETPAYMENT PaymentStatus = "ewallet"
	QRIS
)

type PaymentTransactions struct {
	ID                 uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID             uuid.UUID       `json:"user_id"`
	PaymentType        string          `json:"payment_type"`
	PlanID             uuid.UUID       `json:"plan_id"`
	DurationID         uuid.UUID       `json:"duration_id"`
	TokenQuantity      int             `json:"token_quantity"`
	GrossAmount        decimal.Decimal `json:"gross_amount"`
	DiscountAmount     decimal.Decimal `json:"discount_amount"` // kenapa tidak integer
	FinalAmount        decimal.Decimal `json:"final_amount"`
	PromoCode          *string         `json:"promo_code"`
	XenditInvoiceID    string          `json:"xendit_invoice_id"`
	XenditExternalID   string          `json:"xendit_external_id"`
	PaymentMethod      string          `json:"payment_method"`
	PaymentStatus      string          `json:"payment_status"`
	PaidAt             time.Time       `json:"paid_at"`
	ExpiredAt          time.Time       `json:"expired_at"`
	XenditCallbackData []byte          `json:"xendit_callback_data"`

	Timestamp
}

func (p *PaymentTransactions) TableName() string {
	return "payment_transactions"
}
