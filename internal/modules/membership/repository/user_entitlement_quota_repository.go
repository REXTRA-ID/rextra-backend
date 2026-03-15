package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"rextra-backend/internal/entity"
)

type UserEntitlementQuotaRepository interface {
	Create(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error
	InvalidateByMembershipID(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) error
}

type userEntitlementQuotaRepository struct {
	db *gorm.DB
}

func NewUserEntitlementQuotaRepository(db *gorm.DB) UserEntitlementQuotaRepository {
	return &userEntitlementQuotaRepository{db: db}
}

func (r *userEntitlementQuotaRepository) Create(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Create(&quota).Error
}

func (r *userEntitlementQuotaRepository) InvalidateByMembershipID(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) error {
	if tx == nil { tx = r.db }
	now := time.Now().UTC()
	return tx.WithContext(ctx).
		Model(&entity.UserEntitlementQuota{}).
		Where("membership_id = ? AND (cycle_expired_at IS NULL OR cycle_expired_at > ?)", membershipID, now).
		Update("cycle_expired_at", now).Error
}
