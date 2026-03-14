package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenUsageHistory struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	MembershipID *uuid.UUID `json:"membership_id,omitempty" gorm:"type:uuid;index"`

	EntitlementKey string `json:"entitlement_key" gorm:"type:varchar(200);not null;index"`

	EntitlementName string `json:"entitlement_name" gorm:"type:varchar(200);not null"`

	TokenCost int `json:"token_cost" gorm:"not null;default:1"`

	BalanceBefore int64 `json:"balance_before" gorm:"not null"`
	BalanceAfter  int64 `json:"balance_after" gorm:"not null"`

	ReferenceID *string `json:"reference_id,omitempty" gorm:"type:varchar(500)"`

	UsageLogID *uuid.UUID `json:"usage_log_id,omitempty" gorm:"type:uuid"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index"`

	Membership *Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (TokenUsageHistory) TableName() string {
	return "token_usage_histories"
}

func (t *TokenUsageHistory) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now().UTC()
	return nil
}

func NewTokenUsageHistory(
	userID uuid.UUID,
	membershipID *uuid.UUID,
	entitlementKey string,
	entitlementName string,
	tokenCost int,
	balanceBefore int64,
	balanceAfter int64,
	referenceID *string,
	usageLogID *uuid.UUID,
) TokenUsageHistory {
	return TokenUsageHistory{
		UserID:          userID,
		MembershipID:    membershipID,
		EntitlementKey:  entitlementKey,
		EntitlementName: entitlementName,
		TokenCost:       tokenCost,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		ReferenceID:     referenceID,
		UsageLogID:      usageLogID,
	}
}
