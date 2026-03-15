package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TokenWalletRepository interface {
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.TokenWallet, error)
		CreateWallet(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error
		UpdateBalance(ctx context.Context, tx *gorm.DB, userID string, newBalance int64, ledgerID uuid.UUID) error
	}

	tokenWalletRepository struct {
		db *gorm.DB
	}
)

func NewTokenWalletRepository(db *gorm.DB) TokenWalletRepository {
	return &tokenWalletRepository{db: db}
}

func (r *tokenWalletRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.TokenWallet, error) {
	if tx == nil {
		tx = r.db
	}

	var wallet entity.TokenWallet
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return entity.TokenWallet{}, err
	}

	return wallet, nil
}

func (r *tokenWalletRepository) CreateWallet(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	wallet := entity.TokenWallet{
		UserID:  userID,
		Balance: 0,
	}
	if err := tx.WithContext(ctx).Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

func (r *tokenWalletRepository) UpdateBalance(ctx context.Context, tx *gorm.DB, userID string, newBalance int64, ledgerID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).
		Model(&entity.TokenWallet{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"balance":        newBalance,
			"updated_at":     gorm.Expr("NOW()"),
			"last_ledger_id": ledgerID,
		}).Error
}
