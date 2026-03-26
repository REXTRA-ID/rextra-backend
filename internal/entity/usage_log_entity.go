package entity

import (
	"time"

	"github.com/google/uuid"
)

type UsageLogResult string

const (
	UsageLogResultSuccess UsageLogResult = "SUCCESS"
	UsageLogResultDenied  UsageLogResult = "DENIED"
)

type UsageLog struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipID    uuid.UUID      `gorm:"type:uuid;not null"`
	EntitlementKey  string         `gorm:"type:varchar;not null"`
	EntitlementName string         `gorm:"type:varchar"`
	Result          UsageLogResult `gorm:"type:varchar;not null"`
	ErrorMessage    *string        `gorm:"type:text"`
	ReferenceID     *string        `gorm:"type:varchar"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
}

func NewUsageLog(
	membershipID uuid.UUID,
	entitlementKey string,
	entitlementName string,
	result UsageLogResult,
	errorMessage *string,
	referenceID *string,
) UsageLog {
	return UsageLog{
		MembershipID:    membershipID,
		EntitlementKey:  entitlementKey,
		EntitlementName: entitlementName,
		Result:          result,
		ErrorMessage:    errorMessage,
		ReferenceID:     referenceID,
	}
}
