package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionCycleWriter interface {
	Create(ctx context.Context, tx *gorm.DB, model entity.SubscriptionCycle) (entity.SubscriptionCycle, error)
	GetLastCycleNumber(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) (int, error)
}

type MembershipWriter interface {
	GetByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error)
	Create(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error)
	Update(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error)
}

type PlanDurationReader interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PlanDuration, error)
	GetByPlanAndMonths(ctx context.Context, tx *gorm.DB, planID uuid.UUID, months int) (entity.PlanDuration, error)
}

type MembershipPlanReader interface {
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error)
	GetByPlanName(ctx context.Context, tx *gorm.DB, planName entity.PlanName) (entity.MembershipPlans, error)
}

type DiscountReader interface {
	GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.Discounts, error)
	IncrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	DecrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}

type DiscountRedemptionWriter interface {
	Create(ctx context.Context, tx *gorm.DB, model entity.DiscountRedemption) (entity.DiscountRedemption, error)
	CountByUserAndDiscount(ctx context.Context, tx *gorm.DB, userID, discountID uuid.UUID) (int64, error)
}

type DiscountRedemptionReverser interface {
	Reverse(ctx context.Context, tx *gorm.DB, transactionID string, reason string) error
	GetByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.DiscountRedemption, error)
}

type PoinTransactionWriter interface {
	Create(ctx context.Context, tx *gorm.DB, model entity.PoinTransactions) (entity.PoinTransactions, error)
}

type TokenLedgerWriter interface {
	Create(ctx context.Context, tx *gorm.DB, ledger entity.TokenLedger) (entity.TokenLedger, error)
}

type TokenWalletUpdater interface {
	AddBalance(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int64) (entity.TokenWallet, error)
	GetOrCreateByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.TokenWallet, error)
}

type TopupTransactionWriter interface {
	Create(ctx context.Context, tx *gorm.DB, model *entity.TopupTransaction) (*entity.TopupTransaction, error)
}

type DurationAccessMappingReader interface {
	GetByDurationIDWithEntitlement(ctx context.Context, tx *gorm.DB, durationID uuid.UUID) ([]entity.DurationAccessMapping, error)
}

type UserEntitlementQuotaWriter interface {
	Create(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error
	InvalidateByMembershipID(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) error
}
