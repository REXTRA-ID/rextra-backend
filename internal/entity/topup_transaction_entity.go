package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopupType string

const (
	TopupTypeBundle TopupType = "BUNDLE"
	TopupTypeCustom TopupType = "CUSTOM"
)

type TopupStatus string

const (
	TopupStatusPending  TopupStatus = "PENDING"
	TopupStatusSuccess  TopupStatus = "SUCCESS"
	TopupStatusFailed   TopupStatus = "FAILED"
	TopupStatusExpired  TopupStatus = "EXPIRED"
	TopupStatusRefunded TopupStatus = "REFUNDED"
)

// TopupMetadata for JSONB field
type TopupMetadata map[string]any

// Scan implements sql.Scanner
func (m *TopupMetadata) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

// Value implements driver.Valuer
func (m TopupMetadata) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

type TopupTransaction struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID          uuid.UUID     `json:"user_id" gorm:"not null;column:user_id;index:idx_ledger_user_occurred"`
	Type            TopupType     `json:"type" gorm:"not null;column:type" `
	BundlePackageID *uuid.UUID    `json:"bundle_package_id,omitempty" gorm:"column:bundle_package_id" `
	TokenAmount     int64         `json:"token_amount" gorm:"not null;column:token_amount" `
	TotalPriceRp    int64         `json:"total_price_rp" gorm:"not null;column:total_price_rp" `
	Status          TopupStatus   `json:"status" gorm:"not null;default:'PENDING';column:status;index:idx_topup_status_created" `
	InvoiceID       string        `json:"invoice_id" gorm:"uniqueIndex;column:invoice_id"`
	Provider        *string       `json:"provider,omitempty" gorm:"size:30;column:provider"`
	PaidAt          *time.Time    `json:"paid_at,omitempty" gorm:"column:paid_at" `
	ExpiredAt       *time.Time    `json:"expired_at,omitempty" gorm:"column:expired_at" `
	LedgerID        *uuid.UUID    `json:"ledger_id,omitempty" gorm:"type:uuid;column:ledger_id"`
	Metadata        TopupMetadata `json:"metadata" gorm:"type:jsonb;not null;default:'{}';column:metadata" `
	Timestamp

	User          User                `gorm:"foreignKey:UserID;references:ID" json:"user"`
	BundlePackage *TokenBundlePackage `gorm:"foreignKey:BundlePackageID;references:ID" json:"bundle_package,omitempty"`
	Ledger        *TokenLedger        `gorm:"foreignKey:LedgerID;references:ID" json:"ledger,omitempty"`
}

func (TopupTransaction) TableName() string {
	return "topup_transaction"
}

func (t *TopupTransaction) BeforeSave(tx *gorm.DB) error {
	// Validate type consistency
	if t.Type == TopupTypeBundle && t.BundlePackageID == nil {
		return errors.New("bundle_package_id is required for BUNDLE type")
	}

	if t.Type == TopupTypeCustom && t.BundlePackageID != nil {
		return errors.New("bundle_package_id should be null for CUSTOM type")
	}

	// Validate amounts
	if t.TokenAmount < 1 {
		return errors.New("token_amount must be at least 1")
	}

	if t.TotalPriceRp < 1 {
		return errors.New("total_price_rp must be at least 1")
	}

	return nil
}
