package entity

import (
	"time"

	"github.com/google/uuid"
)

type PromoCodeUsage struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`

	PromoCodeID uuid.UUID  `json:"promo_code_id"`
	Promo       PromoCodes `gorm:"foreignKey:PromoCodeID;references:ID"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	PaymentTransactionID uuid.UUID           `json:"payment_transaction_id"`
	PaymentTransaction   PaymentTransactions `gorm:"foreignKey:PaymentTransactionID;references:ID"`

	UsedAt time.Time `json:"used_at"`
}

func (p *PromoCodeUsage) TableName() string {
	return "promo_code_usage"
}

func (p *PromoCodeUsage) SetPaymentTransactionID(transactionId uuid.UUID) {
	p.PaymentTransactionID = transactionId
}

func NewPromoCodeUsage(promoCodeID uuid.UUID, userID uuid.UUID) PromoCodeUsage {
	return PromoCodeUsage{
		PromoCodeID: promoCodeID,
		UserID:      userID,
		UsedAt:      time.Now().UTC(),
	}
}
