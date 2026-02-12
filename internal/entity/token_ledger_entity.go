package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenDirection string
type TokenSourceType string

const (
	DirectionIN  TokenDirection = "IN"
	DirectionOUT TokenDirection = "OUT"
)

const (
	SourceTypeTopup      TokenSourceType = "TOPUP"
	SourceTypeMembership TokenSourceType = "MEMBERSHIP"
	SourceTypeUsage      TokenSourceType = "USAGE"
	SourceTypeAdjustment TokenSourceType = "ADJUSTMENT"
	SourceTypeRefund     TokenSourceType = "REFUND"
	SourceTypeExpired    TokenSourceType = "EXPIRED"
)

// IsValid checks if the token source type is valid
func (t TokenSourceType) IsValid() bool {
	switch t {
	case SourceTypeTopup, SourceTypeMembership, SourceTypeUsage, SourceTypeAdjustment, SourceTypeRefund, SourceTypeExpired:
		return true
	}
	return false
}

// TokenLedgerMetadata for JSONB field
type TokenLedgerMetadata map[string]any

// Scan implements sql.Scanner
func (m *TokenLedgerMetadata) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &m)
}

// Value implements driver.Valuer
func (m TokenLedgerMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type TokenLedger struct {
	ID            uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OccurredAt    time.Time           `json:"occurred_at" gorm:"not null;default:CURRENT_TIMESTAMP;column:occurred_at;index:idx_ledger_occurred_at"`
	WalletID      uuid.UUID           `json:"wallet_id" gorm:"column:wallet_id;not null"`
	Direction     TokenDirection      `json:"direction" gorm:"varchar(10);not null;column:direction"`
	Amount        int64               `json:"amount" gorm:"not null;column:amount"`
	BalanceBefore int64               `json:"balance_before" gorm:"not null;column:balance_before"`
	BalanceAfter  int64               `json:"balance_after" gorm:"not null;column:balance_after"`
	SourceType    TokenSourceType     `json:"source_type" gorm:"varchar(12);not null;column:source_type"`
	SourceID      *uint64             `json:"source_id" gorm:"column:source_id"`
	ReferenceID   *uuid.UUID          `json:"reference_id" gorm:"column:reference_id"`
	Description   string              `json:"description" gorm:"not null;column:description"`
	Metadata      TokenLedgerMetadata `json:"metadata" gorm:"type:jsonb;not null;default:'{}';column:metadata"`
	OperatorID    *uuid.UUID          `json:"operator_id,omitempty" gorm:"column:operator_id"`
	CreatedAt     time.Time           `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at"`

	// Relations
	Wallet   TokenWallet `gorm:"foreignKey:WalletID;references:ID"`
	Operator *User       `gorm:"foreignKey:OperatorID;references:ID"`
	Ref      *uuid.UUID  `gorm:"foreignKey:ReferenceID;references:ID"`
}

func (TokenLedger) TableName() string {
	return "token_ledger"
}

func (l *TokenLedger) BeforeCreate(tx *gorm.DB) error {
	if l.Amount <= 0 {
		return errors.New("ledger: amount must be positive")
	}

	if err := l.validateBalanceConsistency(); err != nil {
		return err
	}

	return nil
}

func (l *TokenLedger) validateBalanceConsistency() error {
	var expectedBalance int64

	switch l.Direction {
	case DirectionIN:
		expectedBalance = l.BalanceBefore + l.Amount
	case DirectionOUT:
		expectedBalance = l.BalanceBefore - l.Amount
	default:
		return errors.New("ledger: invalid direction")
	}

	if l.BalanceAfter != expectedBalance {
		return fmt.Errorf("ledger: balance mismatch, expected %d but got %d", expectedBalance, l.BalanceAfter)
	}

	return nil
}
