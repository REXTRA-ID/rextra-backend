package service

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanRepository interface {
	GetByIDWithDurations(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error)
}

type PromoRepository interface {
	CountEligibleVouchers(ctx context.Context, planID string) (int, error)
	ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, subtotal int64) (int64, error)
	RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID string) error
}

type UserRepository interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.User, error)
}
