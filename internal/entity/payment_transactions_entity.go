package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentType string
type ChangeType string
type PaymentProvider string
type CancelReason string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusCanceled PaymentStatus = "cancelled"

	PaymentTypeMembership PaymentType = "membership"

	ChangeTypeNewPurchase ChangeType = "PEMBELIAN_BARU"
	ChangeTypeRenewal     ChangeType = "RENEWAL"
	ChangeTypeUpgrade     ChangeType = "UPGRADE"
	ChangeTypeDowngrade   ChangeType = "DOWNGRADE"

	PaymentProviderTripay PaymentProvider = "TRIPAY"

	CancelReasonChangeOrder  CancelReason = "CHANGE_ORDER"
	CancelReasonBudget       CancelReason = "BUDGET"
	CancelReasonNotNow       CancelReason = "NOT_NOW"
	CancelReasonPaymentIssue CancelReason = "PAYMENT_ISSUE"
	CancelReasonOther        CancelReason = "OTHER"
)

type PaymentTransactions struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	TransactionID string `json:"transaction_id" gorm:"type:varchar(100);uniqueIndex;not null"`

	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	MembershipID *uuid.UUID `json:"membership_id,omitempty" gorm:"type:uuid;index"`

	ChangeType ChangeType `json:"change_type" gorm:"type:varchar(30);not null"`

	FromPlan           *string `json:"from_plan,omitempty" gorm:"type:varchar(20)"`
	ToPlan             string  `json:"to_plan" gorm:"type:varchar(20);not null"`
	FromDurationMonths *int    `json:"from_duration_months,omitempty"`
	ToDurationMonths   int     `json:"to_duration_months" gorm:"not null;default:1"`

	PlanID     *uuid.UUID `json:"plan_id,omitempty" gorm:"type:uuid;index"`
	DurationID *uuid.UUID `json:"duration_id,omitempty" gorm:"type:uuid;index"`

	TokenBundlePackageID *uuid.UUID `json:"token_bundle_package_id,omitempty" gorm:"type:uuid;index"`
	TokenTopupAmount     int64      `json:"token_topup_amount" gorm:"not null;default:0"`
	TokenTopupPrice      int64      `json:"token_topup_price" gorm:"not null;default:0"`

	SubtotalAmount int64 `json:"subtotal_amount" gorm:"not null;default:0"`
	DiscountAmount int64 `json:"discount_amount" gorm:"not null;default:0"`
	DurationCredit int64 `json:"duration_credit" gorm:"not null;default:0"`
	AdminFee       int64 `json:"admin_fee" gorm:"not null;default:0"`
	TotalAmount    int64 `json:"total_amount" gorm:"not null;default:0"`

	PromoCode *string `json:"promo_code,omitempty" gorm:"type:varchar(50)"`

	PaymentProvider   PaymentProvider `json:"payment_provider" gorm:"type:varchar(20)"`
	PaymentMethod     *string         `json:"payment_method,omitempty" gorm:"type:varchar(100)"`
	PaymentStatus     PaymentStatus   `json:"payment_status" gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentExternalID string          `json:"payment_external_id" gorm:"type:varchar(200)"`
	MerchantRef       string          `json:"merchant_ref" gorm:"type:varchar(200)"`
	PaymentURL        *string         `json:"payment_url,omitempty" gorm:"type:varchar(500)"`
	PayCode           *string         `json:"pay_code,omitempty" gorm:"type:varchar(100)"`

	UserName  string `json:"user_name" gorm:"type:varchar(200)"`
	UserEmail string `json:"user_email" gorm:"type:varchar(200)"`

	Items datatypes.JSON `json:"items" gorm:"type:jsonb;default:'[]'"`

	PrimaryStatus string `json:"primary_status" gorm:"type:varchar(50)"`

	PaidAt     *time.Time `json:"paid_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`

	CancelReason *CancelReason `json:"cancel_reason,omitempty" gorm:"type:varchar(30)"`
	CancelNote   *string       `json:"cancel_note,omitempty" gorm:"type:text"`

	PaymentCallbackData datatypes.JSON `json:"payment_callback_data,omitempty" gorm:"type:jsonb"`

	Timestamp

	Plan               *MembershipPlans    `json:"plan,omitempty" gorm:"foreignKey:PlanID;references:ID"`
	Duration           *PlanDuration       `json:"duration,omitempty" gorm:"foreignKey:DurationID;references:ID"`
	Membership         *Memberships        `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
	TokenBundlePackage *TokenBundlePackage `json:"token_bundle_package,omitempty" gorm:"foreignKey:TokenBundlePackageID;references:ID"`
}

func (PaymentTransactions) TableName() string {
	return "payment_transactions"
}

func (p *PaymentTransactions) BeforeCreate(tx *gorm.DB) error {
	if p.TotalAmount < 0 {
		return errors.New("total_amount tidak boleh negatif")
	}
	return nil
}

func NewPaymentTransaction(
	transactionID string,
	userID uuid.UUID,
	membershipID *uuid.UUID,
	changeType ChangeType,
	fromPlan *string,
	toPlan string,
	fromDurationMonths *int,
	toDurationMonths int,
	planID, durationID *uuid.UUID,
	tokenBundlePackageID *uuid.UUID,
	tokenTopupAmount int64,
	tokenTopupPrice int64,
	subtotalAmount, discountAmount, durationCredit, adminFee, totalAmount int64,
	promoCode *string,
	paymentMethod string,
	merchantRef string,
	userName, userEmail string,
	items datatypes.JSON,
	expiresAt *time.Time,
) PaymentTransactions {
	method := paymentMethod
	return PaymentTransactions{
		TransactionID:        transactionID,
		UserID:               userID,
		MembershipID:         membershipID,
		ChangeType:           changeType,
		FromPlan:             fromPlan,
		ToPlan:               toPlan,
		FromDurationMonths:   fromDurationMonths,
		ToDurationMonths:     toDurationMonths,
		PlanID:               planID,
		DurationID:           durationID,
		TokenBundlePackageID: tokenBundlePackageID,
		TokenTopupAmount:     tokenTopupAmount,
		TokenTopupPrice:      tokenTopupPrice,
		SubtotalAmount:       subtotalAmount,
		DiscountAmount:       discountAmount,
		DurationCredit:       durationCredit,
		AdminFee:             adminFee,
		TotalAmount:          totalAmount,
		PromoCode:            promoCode,
		PaymentProvider:      PaymentProviderTripay,
		PaymentMethod:        &method,
		PaymentStatus:        PaymentStatusPending,
		MerchantRef:          merchantRef,
		PrimaryStatus:        "BERLANGSUNG",
		UserName:             userName,
		UserEmail:            userEmail,
		Items:                items,
		ExpiresAt:            expiresAt,
	}
}

func (p *PaymentTransactions) SetTripayResponse(reference, checkoutURL, payCode string, expiredTime int64) {
	p.PaymentExternalID = reference
	p.PaymentURL = &checkoutURL
	if payCode != "" {
		p.PayCode = &payCode
	}
	if expiredTime > 0 {
		t := time.Unix(expiredTime, 0).UTC()
		p.ExpiresAt = &t
	}
}

func (p *PaymentTransactions) MarkPaid(paidAt time.Time, callbackData []byte) {
	p.PaymentStatus = PaymentStatusPaid
	p.PaidAt = &paidAt
	p.PrimaryStatus = "BERHASIL"
	p.PaymentCallbackData = datatypes.JSON(callbackData)
	p.UpdatedAt = time.Now().UTC()
}

func (p *PaymentTransactions) MarkFailed(callbackData []byte) {
	p.PaymentStatus = PaymentStatusFailed
	p.PrimaryStatus = "GAGAL"
	p.PaymentCallbackData = datatypes.JSON(callbackData)
	p.UpdatedAt = time.Now().UTC()
}

func (p *PaymentTransactions) MarkExpired(callbackData []byte) {
	p.PaymentStatus = PaymentStatusExpired
	p.PrimaryStatus = "KEDALUWARSA"
	p.PaymentCallbackData = datatypes.JSON(callbackData)
	p.UpdatedAt = time.Now().UTC()
}

func (p *PaymentTransactions) Cancel(reason CancelReason, note *string) error {
	if p.PaymentStatus != PaymentStatusPending {
		return errors.New("only PENDING transactions can be cancelled")
	}
	now := time.Now().UTC()
	p.PaymentStatus = PaymentStatusCanceled
	p.CanceledAt = &now
	p.CancelReason = &reason
	p.CancelNote = note
	p.PrimaryStatus = "DIBATALKAN"
	p.UpdatedAt = now
	return nil
}

func (p *PaymentTransactions) IsPaid() bool {
	return p.PaymentStatus == PaymentStatusPaid
}

func (p *PaymentTransactions) IsCancellableByUser() bool {
	return p.PaymentStatus == PaymentStatusPending
}
