package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type EnumFeature string

const (
	KENALIDIRI    EnumFeature = "kenali_diri"
	CVGENERATOR   EnumFeature = "cv_generator"
	AIINTERVIEWER EnumFeature = "ai_interviewer"
)

type TokenUsageHistory struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID uuid.UUID `json:"user_id"`

	TokenTransaction   TokenTransaction `gorm:"foreignKey:TokenTransactionID;references:ID"`
	TokenTransactionID uuid.UUID        `json:"token_transaction_id"`

	FeatureName   EnumFeature    `json:"feature_name" gorm:"varchar(100)"`
	TokenUsed     int            `json:"token_used"`
	UsageMetadata datatypes.JSON `json:"usage_metadata"`

	CreatedAt time.Time `json:"created_at"`
}

func (t *TokenUsageHistory) TableName() string {
	return "token_usage_history"
}

func NewTokenUsageHistory(userId uuid.UUID, tokenTransactionId uuid.UUID, feature string, metadata []byte) TokenUsageHistory {

	return TokenUsageHistory{
		UserID:             userId,
		TokenTransactionID: tokenTransactionId,
		FeatureName:        EnumFeature(feature),
		UsageMetadata:      datatypes.JSON(metadata),
		CreatedAt:          time.Now().UTC(),
	}
}
