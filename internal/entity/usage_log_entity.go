package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsageLogResult string

const (
	UsageLogResultSuccess UsageLogResult = "success"
	UsageLogResultDenied  UsageLogResult = "denied"
)

type UsageLog struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	MembershipID *uuid.UUID `json:"membership_id,omitempty" gorm:"type:uuid;index"`

	EntitlementKey string `json:"entitlement_key" gorm:"type:varchar(200);not null;index"`

	EntitlementName string `json:"entitlement_name" gorm:"type:varchar(200);not null"`

	Result UsageLogResult `json:"result" gorm:"type:varchar(20);not null;index"`

	ErrorMessage *string `json:"error_message,omitempty" gorm:"type:text"`

	ReferenceID *string `json:"reference_id,omitempty" gorm:"type:varchar(500)"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;index"`

	Membership *Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (UsageLog) TableName() string {
	return "usage_logs"
}

func (u *UsageLog) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now().UTC()
	return nil
}

func NewUsageLog(
	membershipID uuid.UUID,
	entitlementKey string,
	entitlementName string,
	result UsageLogResult,
	errorMessage *string,
	referenceID *string,
) UsageLog {
	var memID *uuid.UUID
	if membershipID != uuid.Nil {
		memID = &membershipID
	}

	return UsageLog{
		MembershipID:    memID,
		EntitlementKey:  entitlementKey,
		EntitlementName: entitlementName,
		Result:          result,
		ErrorMessage:    errorMessage,
		ReferenceID:     referenceID,
	}
}

func (u *UsageLog) IsGranted() bool {
	return u.Result == UsageLogResultSuccess
}
