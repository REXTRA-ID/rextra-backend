package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntitlementQuota struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipID uuid.UUID `json:"membership_id" gorm:"type:uuid;not null;index:idx_ueq_unique,unique"`
	EntitlementID uuid.UUID `json:"entitlement_id" gorm:"type:uuid;not null;index:idx_ueq_unique,unique"`

	EntitlementKey string `json:"entitlement_key" gorm:"type:varchar(200);not null;index"`

	QuotaGranted int `json:"quota_granted" gorm:"not null;default:0"`

	QuotaUsed int `json:"quota_used" gorm:"not null;default:0"`

	QuotaRemaining int `json:"quota_remaining" gorm:"not null;default:0"`

	CycleStartedAt time.Time  `json:"cycle_started_at" gorm:"not null"`
	CycleExpiredAt *time.Time `json:"cycle_expired_at,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	Membership Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	Entitlement Entitlement `json:"entitlement,omitempty" gorm:"foreignKey:EntitlementID;references:ID"`
}

func (UserEntitlementQuota) TableName() string {
	return "user_entitlement_quotas"
}

func (u *UserEntitlementQuota) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

func NewUserEntitlementQuota(
	membershipID uuid.UUID,
	entitlementID uuid.UUID,
	entitlementKey string,
	quotaGranted int,
	cycleStartedAt time.Time,
	cycleExpiredAt *time.Time,
) UserEntitlementQuota {
	return UserEntitlementQuota{
		MembershipID:   membershipID,
		EntitlementID:  entitlementID,
		EntitlementKey: entitlementKey,
		QuotaGranted:   quotaGranted,
		QuotaUsed:      0,
		QuotaRemaining: quotaGranted,
		CycleStartedAt: cycleStartedAt,
		CycleExpiredAt: cycleExpiredAt,
	}
}

func (u *UserEntitlementQuota) ConsumeOne() {
	u.QuotaUsed++
	u.QuotaRemaining--
	u.UpdatedAt = time.Now().UTC()
}

func (u *UserEntitlementQuota) IsExhausted() bool {
	return u.QuotaRemaining <= 0
}

func (u *UserEntitlementQuota) IsActive() bool {
	now := time.Now().UTC()
	if now.Before(u.CycleStartedAt) {
		return false
	}
	if u.CycleExpiredAt != nil && now.After(*u.CycleExpiredAt) {
		return false
	}
	return true
}
