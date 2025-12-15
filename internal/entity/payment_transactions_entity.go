package entity

import (
	"errors"
	"rextra-backend/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type (
	PaymentStatus string

	PaymentType string

	PaymentServiceStatus string
)

const (
	PENDING PaymentStatus = "pending"
	PAID    PaymentStatus = "paid"
	FAILED  PaymentStatus = "failed"
	EXPIRED PaymentStatus = "expired"

	XENDITSUCCEEDED PaymentServiceStatus = "SUCCEEDED"
	XENDITFAILED    PaymentServiceStatus = "FAILED"
	XENDITEXPIRED   PaymentServiceStatus = "EXPIRED"

	MIDTRANSSETTLEMENT PaymentServiceStatus = "settlement"
	MIDTRANSCAPTURE    PaymentServiceStatus = "capture"
	MIDTRANSFAILURE    PaymentServiceStatus = "failure"
	MIDTRANSEXPIRE     PaymentServiceStatus = "failure"

	MEMBERSHIP      PaymentType = "membership"
	TOKENSTANDALONE PaymentType = "token_standalone"
)

type PaymentTransactions struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID     `json:"user_id"`
	PaymentType   PaymentType   `json:"payment_type"`
	PaymentMethod string        `json:"payment_method"`
	PaymentStatus PaymentStatus `json:"payment_status"`

	PlanID *uuid.UUID       `json:"plan_id"`
	Plan   *MembershipPlans `gorm:"foreignKey:PlanID"`

	DurationID *uuid.UUID          `json:"duration_id"`
	Duration   *MembershipDuration `gorm:"foreignKey:DurationID"`

	TokenQuantity int `json:"token_quantity" gorm:"default:0"`

	GrossAmount         float64        `json:"gross_amount"`
	FinalAmount         float64        `json:"final_amount"`
	PromoCode           *string        `json:"promo_code"`
	PaymentInvoiceID    string         `json:"payment_invoice_id"`
	PaymentExternalID   string         `json:"payment_external_id"`
	PaidAt              *time.Time     `json:"paid_at"`
	ExpiredAt           *time.Time     `json:"expired_at"`
	PaymentCallbackData datatypes.JSON `json:"xendit_callback_data"`

	Timestamp
}

func (p *PaymentTransactions) TableName() string {
	return "payment_transactions"
}

func NewPaymenTransaction(userId uuid.UUID,
	paymentType string,
	planId, durationId *uuid.UUID,
	grossAmount float64,
	tokenQuantity int) PaymentTransactions {

	payemntExternalId := utils.PaymentExternalID(paymentType, userId.String())

	return PaymentTransactions{
		UserID:            userId,
		PlanID:            planId,
		DurationID:        durationId,
		TokenQuantity:     tokenQuantity,
		PaymentType:       PaymentType(paymentType),
		PaymentStatus:     PENDING,
		GrossAmount:       grossAmount,
		FinalAmount:       grossAmount,
		PaymentExternalID: payemntExternalId,
	}
}

func (p *PaymentTransactions) SetInvoice(invoiceId string) {
	p.PaymentInvoiceID = invoiceId
}

func (p *PaymentTransactions) UpdateXenditTransaction(paymentStatus string, paymentCallback []byte) bool {
	p.PaymentCallbackData = datatypes.JSON(paymentCallback)
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

func (p *PaymentTransactions) UpdateMidtransTransaction(
	paymentStatus, transactionTime, paymentMethod string,
	paymentCallback []byte) error {

	p.PaymentCallbackData = datatypes.JSON(paymentCallback)
	t, err := time.Parse("2006-01-02 15:04:05", transactionTime)
	if err != nil {
		return err
	}

	if paymentStatus == string(MIDTRANSSETTLEMENT) || paymentStatus == string(MIDTRANSCAPTURE) {
		p.PaidAt = &t
		p.PaymentStatus = PAID
		p.PaymentMethod = paymentMethod
		return nil
	} else if paymentStatus == string(MIDTRANSEXPIRE) {
		p.ExpiredAt = &t
		p.PaymentStatus = EXPIRED
		p.PaymentMethod = ""
		return nil
	} else if paymentStatus == string(MIDTRANSFAILURE) {
		p.PaymentStatus = FAILED
		p.PaymentMethod = ""
		return nil
	} else {
		return errors.New("invalid status")
	}
}
