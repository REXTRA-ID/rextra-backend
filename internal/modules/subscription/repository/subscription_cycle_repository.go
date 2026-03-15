package repository

import (
	"context"
	"rextra-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	SubscriptionCycleRepository interface {
		GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, planName, status string, dateFrom, dateTo time.Time, userID uuid.UUID) ([]entity.SubscriptionCycle, int64, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.SubscriptionCycle, error)
		GetAllByMembershipID(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) ([]entity.SubscriptionCycle, error)
		GetLastCycleNumber(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) (int, error)
		Create(ctx context.Context, tx *gorm.DB, model entity.SubscriptionCycle) (entity.SubscriptionCycle, error)
	}

	subscriptionCycleRepository struct {
		db *gorm.DB
	}
)

func NewSubscriptionCycleRepository(db *gorm.DB) SubscriptionCycleRepository {
	return &subscriptionCycleRepository{db: db}
}

func (r *subscriptionCycleRepository) Create(ctx context.Context, tx *gorm.DB, model entity.SubscriptionCycle) (entity.SubscriptionCycle, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.SubscriptionCycle{}, err }
	return model, nil
}

func (r *subscriptionCycleRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, planName, status string, dateFrom, dateTo time.Time, userID uuid.UUID) ([]entity.SubscriptionCycle, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.SubscriptionCycle{})
	if planName != "" { query = query.Where("subscription_cycles.plan_name = ?", planName) }
	if status != "" { query = query.Where("subscription_cycles.status = ?", status) }
	if !dateFrom.IsZero() { query = query.Where("subscription_cycles.start_date >= ?", dateFrom) }
	if !dateTo.IsZero() { query = query.Where("subscription_cycles.start_date < ?", dateTo.AddDate(0, 0, 1)) }
	if userID != uuid.Nil {
		query = query.Joins("JOIN memberships ON memberships.id = subscription_cycles.membership_id").Where("memberships.user_id = ?", userID)
	}
	var total int64
	if err := query.Select("COUNT(DISTINCT subscription_cycles.id)").Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.SubscriptionCycle
	if err := query.Order("start_date DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *subscriptionCycleRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.SubscriptionCycle, error) {
	if tx == nil { tx = r.db }
	var result entity.SubscriptionCycle
	if err := tx.WithContext(ctx).Preload("Membership").Preload("PlanDuration").First(&result, "id = ?", id).Error; err != nil { return entity.SubscriptionCycle{}, err }
	return result, nil
}

func (r *subscriptionCycleRepository) GetAllByMembershipID(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) ([]entity.SubscriptionCycle, error) {
	if tx == nil { tx = r.db }
	var results []entity.SubscriptionCycle
	if err := tx.WithContext(ctx).Where("membership_id = ?", membershipID).Order("cycle_number DESC").Find(&results).Error; err != nil { return nil, err }
	return results, nil
}

func (r *subscriptionCycleRepository) GetLastCycleNumber(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID) (int, error) {
	if tx == nil { tx = r.db }
	var maxCycle int
	err := tx.WithContext(ctx).Raw("SELECT COALESCE(MAX(cycle_number), 0) FROM subscription_cycles WHERE membership_id = ?", membershipID).Scan(&maxCycle).Error
	return maxCycle, err
}
