package entity

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type PaymentTransactions struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	PaymentType        string
	PlanID             uuid.UUID
	DurationID         uuid.UUID
	TokenQuantity      int
	GrossAmount        decimal.Decimal
	DiscountAmount     decimal.Decimal
	FinalAmount        decimal.Decimal
	PromoCode          *string
	XenditInvoiceID    string
	XenditExternalID   string
	PaymentMethod      string
	PaymentStatus      string
	PaidAt             time.Time
	ExpiredAt          time.Time
	XenditCallbackData []byte

	Timestamp
}
