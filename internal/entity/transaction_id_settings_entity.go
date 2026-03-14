package entity

import (
	"time"

	"github.com/google/uuid"
)

type TrxPattern string

const (
	TrxPatternDateDaily   TrxPattern = "DATE_DAILY"
	TrxPatternDateMonthly TrxPattern = "DATE_MONTHLY"
	TrxPatternSequential  TrxPattern = "SEQUENTIAL"
)

type TransactionIdSettings struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	TrxPrefix  string     `json:"trx_prefix" gorm:"type:varchar(20);not null;default:'TRX'"`
	TrxPattern TrxPattern `json:"trx_pattern" gorm:"type:varchar(30);not null;default:'DATE_DAILY'"`

	UpdatedBy string    `json:"updated_by" gorm:"type:varchar(100);default:'System'"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (TransactionIdSettings) TableName() string {
	return "transaction_id_settings"
}

func DefaultTransactionIdSettings() TransactionIdSettings {
	return TransactionIdSettings{
		ID:         uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		TrxPrefix:  "TRX",
		TrxPattern: TrxPatternDateDaily,
		UpdatedBy:  "System",
		UpdatedAt:  time.Now().UTC(),
	}
}
