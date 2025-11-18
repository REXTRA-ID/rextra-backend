package entity

import (
	"time"

	"rextra-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type (
	PaymentStatus string

	PaymentType string

	XenditPaymentStatus string
)

const (
	PENDING PaymentStatus = "pending"
	PAID    PaymentStatus = "paid"
	FAILED  PaymentStatus = "failed"
	EXPIRED PaymentStatus = "expired"

	XENDITSUCCEEDED XenditPaymentStatus = "SUCCEEDED"
	XENDITFAILED    XenditPaymentStatus = "FAILED"
	XENDITEXPIRED   XenditPaymentStatus = "EXPIRED"

	MEMBERSHIP     PaymentType   = "membership"
	TOKENSTNDALONE PaymentStatus = "token_standalone"
)

type PaymentTransactions struct {
	ID            uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID   `json:"user_id"`
	PaymentType   PaymentType `json:"payment_type"`
	PlanID        *uuid.UUID  `json:"plan_id"`
	DurationID    *uuid.UUID  `json:"duration_id"`
	TokenQuantity int         `json:"token_quantity" gorm:"default:0"`
	GrossAmount   float64     `json:"gross_amount"`
	// DiscountAmount     float64 `json:"discount_amount"` discount untuk sekarang tidak dibutuhkan
	FinalAmount        float64        `json:"final_amount"`
	PromoCode          *string        `json:"promo_code"`
	XenditInvoiceID    string         `json:"xendit_invoice_id"`
	XenditExternalID   string         `json:"xendit_external_id"`
	PaymentMethod      string         `json:"payment_method"`
	PaymentStatus      PaymentStatus  `json:"payment_status"`
	PaidAt             *time.Time     `json:"paid_at"`
	ExpiredAt          *time.Time     `json:"expired_at"`
	XenditCallbackData datatypes.JSON `json:"xendit_callback_data"`

	Timestamp
}

func (p *PaymentTransactions) TableName() string {
	return "payment_transactions"
}

func (p *PaymentTransactions) CalculatePrice(plan MembershipPlans, duration MembershipDuration) float64 {
	price := plan.BaseMonthlyPrice * float64(duration.DurationMonth)
	return price
}

func NewPaymenTransaction(userId uuid.UUID,
	paymentType, PaymentMethod string,
	planId, durationId *uuid.UUID,
	grossAmount float64,
	xenditInvoiceId string,
	tokenQuantity int) PaymentTransactions {

	xenditExternalId := utils.MakeXenditExternalID(paymentType, userId.String())

	return PaymentTransactions{
		UserID:           userId,
		PlanID:           planId,
		DurationID:       durationId,
		TokenQuantity:    tokenQuantity,
		PaymentType:      PaymentType(paymentType),
		PaymentStatus:    PENDING,
		PaymentMethod:    PaymentMethod,
		GrossAmount:      grossAmount,
		FinalAmount:      grossAmount,
		XenditInvoiceID:  xenditInvoiceId,
		XenditExternalID: xenditExternalId,
	}
}

func (p *PaymentTransactions) UpdateTransaction(paymentStatus string, xenditCallback []byte) bool {
	p.XenditCallbackData = datatypes.JSON(xenditCallback)
	now := time.Now().UTC()

	if paymentStatus == string(XENDITEXPIRED) {
		p.PaymentStatus = EXPIRED
		p.ExpiredAt = &now
		return false
	} else if paymentStatus == string(XENDITFAILED) {
		p.PaymentStatus = FAILED
		return false
	}

	p.PaymentStatus = PAID
	p.PaidAt = &now
	return true
}
