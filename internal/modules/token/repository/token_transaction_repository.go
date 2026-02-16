package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, tokenTransaction entity.TokenTransaction) (entity.TokenTransaction, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.TokenTransaction, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenTransaction, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.TokenTransaction, error)
	}

	tokenTransactionRepository struct {
		db *gorm.DB
	}
)

func NewTokenTransactionRepository(db *gorm.DB) TokenTransactionRepository {
	return &tokenTransactionRepository{
		db: db,
	}
}

func (r *tokenTransactionRepository) Create(ctx context.Context, tx *gorm.DB, tokenTransaction entity.TokenTransaction) (entity.TokenTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&tokenTransaction).Error; err != nil {
		return entity.TokenTransaction{}, err
	}

	return tokenTransaction, nil
}

func (r *tokenTransactionRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.TokenTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var userTokenTransaction []entity.TokenTransaction
	if err := tx.WithContext(ctx).Find(&userTokenTransaction, "user_id = ?", userId).Error; err != nil {
		return []entity.TokenTransaction{}, err
	}

	return userTokenTransaction, nil
}

func (r *tokenTransactionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var tokenTransactions []entity.TokenTransaction
	if err := tx.WithContext(ctx).Find(&tokenTransactions).Error; err != nil {
		return []entity.TokenTransaction{}, err
	}

	return tokenTransactions, nil
}

func (r *tokenTransactionRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.TokenTransaction, error) {
	return []entity.TokenTransaction{}, nil
}
