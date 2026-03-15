package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	UserEntitlementRepository interface {
		GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error)
		GetEntitlementsByPlanDurationID(ctx context.Context, planDurationID *uuid.UUID) ([]entity.DurationAccessMapping, error)
		GetQuotasByMembershipID(ctx context.Context, membershipID uuid.UUID) ([]entity.UserEntitlementQuota, error)
	}

	userEntitlementRepository struct {
		db *gorm.DB
	}
)

func NewUserEntitlementRepository(db *gorm.DB) UserEntitlementRepository {
	return &userEntitlementRepository{db: db}
}

func (r *userEntitlementRepository) GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error) {
	var membership entity.Memberships
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true AND (expired_at > NOW() OR expired_at IS NULL)", userID).First(&membership).Error
	return membership, err
}

func (r *userEntitlementRepository) GetEntitlementsByPlanDurationID(ctx context.Context, planDurationID *uuid.UUID) ([]entity.DurationAccessMapping, error) {
	var mappings []entity.DurationAccessMapping
	err := r.db.WithContext(ctx).Preload("Entitlement").Joins("JOIN entitlements ON entitlements.id = duration_access_mappings.entitlement_id").Where("duration_access_mappings.plan_duration_id = ? AND duration_access_mappings.status = 'aktif' AND entitlements.status = 'active'", *planDurationID).Find(&mappings).Error
	return mappings, err
}

func (r *userEntitlementRepository) GetQuotasByMembershipID(ctx context.Context, membershipID uuid.UUID) ([]entity.UserEntitlementQuota, error) {
	var quotas []entity.UserEntitlementQuota
	err := r.db.WithContext(ctx).Where("membership_id = ?", membershipID).Find(&quotas).Error
	return quotas, err
}
