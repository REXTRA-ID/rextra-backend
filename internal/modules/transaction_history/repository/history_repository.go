package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	HistoryRepository interface {
		GetMembershipTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.PaymentTransactions, int64, error)
		GetTopupTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.TopupTransaction, int64, error)
	}

	historyRepository struct {
		db *gorm.DB
	}
)

func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) GetMembershipTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.PaymentTransactions, int64, error) {
	var records []entity.PaymentTransactions
	var total int64
	base := r.db.WithContext(ctx).Model(&entity.PaymentTransactions{}).Where("user_id = ?", userID)
	if err := base.Count(&total).Error; err != nil { return nil, 0, err }
	if err := base.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil { return nil, 0, err }
	return records, total, nil
}

func (r *historyRepository) GetTopupTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.TopupTransaction, int64, error) {
	var records []entity.TopupTransaction
	var total int64
	base := r.db.WithContext(ctx).Model(&entity.TopupTransaction{}).Where("user_id = ?", userID)
	if err := base.Count(&total).Error; err != nil { return nil, 0, err }
	if err := base.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil { return nil, 0, err }
	return records, total, nil
}
