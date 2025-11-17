package entity

import (
	"time"

	"github.com/google/uuid"
)

type TokenUsageHistory struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID             uuid.UUID `json:"user_id"`
	TokenTransactionID uuid.UUID `json:"token_transaction_id"` // seperti nya tidak diperlukan
	FeatureName        string    `json:"feature_name"`
	TokenUsed          int       `json:"token_used"`
	UsageMetadata      []byte    `json:"usage_metadata"` // masih perlu bertanya

	CreatedAt time.Time `json:"created_at"`
}

func (t *TokenUsageHistory) TableName() string {
	return "token_usage_history"
}
