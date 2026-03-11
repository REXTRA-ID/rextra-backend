package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenUsageHistoryRepository interface {
		Create(ctx context.Context, tx *gorm.DB, tokenUsage entity.TokenUsageHistory) (entity.TokenUsageHistory, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.TokenUsageHistory, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenUsageHistory, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.TokenUsageHistory, error)
	}

	tokenUsageHistoryRepository struct {
		db *gorm.DB
	}
)

func NewTokenUsageHistoryRepository(db *gorm.DB) TokenUsageHistoryRepository {
	return &tokenUsageHistoryRepository{
		db: db,
	}
}

func (r *tokenUsageHistoryRepository) Create(ctx context.Context, tx *gorm.DB, tokenUsage entity.TokenUsageHistory) (entity.TokenUsageHistory, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&tokenUsage).Error; err != nil {
		return entity.TokenUsageHistory{}, err
	}

	return tokenUsage, nil
}

func (r *tokenUsageHistoryRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.TokenUsageHistory, error) {
	if tx == nil {
		tx = r.db
	}

	var usersTokenUsage []entity.TokenUsageHistory
	if err := tx.WithContext(ctx).Find(&usersTokenUsage, "user_id = ?", userId).Error; err != nil {
		return []entity.TokenUsageHistory{}, err
	}

	return usersTokenUsage, nil
}

func (r *tokenUsageHistoryRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenUsageHistory, error) {
	return []entity.TokenUsageHistory{}, nil
}

func (r *tokenUsageHistoryRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.TokenUsageHistory, error) {
	return []entity.TokenUsageHistory{}, nil
}
