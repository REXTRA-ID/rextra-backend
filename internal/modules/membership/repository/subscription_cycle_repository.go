// package repository

package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	SubscriptionCycleRepository interface {
		GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entity.SubscriptionCycle, error)
		GetTotalCount(ctx context.Context, tx *gorm.DB) (int64, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId string) ([]entity.SubscriptionCycle, error)
		Create(ctx context.Context, tx *gorm.DB, data entity.SubscriptionCycle) (entity.SubscriptionCycle, error)
	}

	subscriptionCycleRepository struct {
		db *gorm.DB
	}
)

func NewSubscriptionCycleRepository(db *gorm.DB) SubscriptionCycleRepository {
	return &subscriptionCycleRepository{db: db}
}

func (r *subscriptionCycleRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int) ([]entity.SubscriptionCycle, error) {
	if tx == nil {
		tx = r.db
	}

	var results []entity.SubscriptionCycle
	if err := tx.WithContext(ctx).
		Model(&entity.SubscriptionCycle{}).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *subscriptionCycleRepository) GetTotalCount(ctx context.Context, tx *gorm.DB) (int64, error) {
	if tx == nil {
		tx = r.db
	}

	var total int64
	if err := tx.WithContext(ctx).Model(&entity.SubscriptionCycle{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *subscriptionCycleRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId string) ([]entity.SubscriptionCycle, error) {
	if tx == nil {
		tx = r.db
	}

	var userSubs []entity.SubscriptionCycle
	if err := tx.WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at desc").
		Find(&userSubs).Error; err != nil {
		return nil, err
	}

	return userSubs, nil
}

func (r *subscriptionCycleRepository) Create(ctx context.Context, tx *gorm.DB, data entity.SubscriptionCycle) (entity.SubscriptionCycle, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.SubscriptionCycle{}, err
	}

	return data, nil
}
