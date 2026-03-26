package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PoinTransactionsRepository interface {
		Create(ctx context.Context, tx *gorm.DB, poinTransaction entity.PoinTransactions) (entity.PoinTransactions, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.PoinTransactions, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.PoinTransactions, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.PoinTransactions, error)
	}

	poinTransactionRepository struct {
		db *gorm.DB
	}
)

func NewPoinTransactionsRepository(db *gorm.DB) PoinTransactionsRepository {
	return &poinTransactionRepository{
		db: db,
	}
}

func (r *poinTransactionRepository) Create(ctx context.Context, tx *gorm.DB, poinTransaction entity.PoinTransactions) (entity.PoinTransactions, error) {

	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&poinTransaction).Error; err != nil {
		return entity.PoinTransactions{}, err
	}

	return poinTransaction, nil
}

func (r *poinTransactionRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.PoinTransactions, error) {
	return []entity.PoinTransactions{}, nil
}

func (r *poinTransactionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.PoinTransactions, error) {
	return []entity.PoinTransactions{}, nil
}

func (r *poinTransactionRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.PoinTransactions, error) {
	return []entity.PoinTransactions{}, nil
}
