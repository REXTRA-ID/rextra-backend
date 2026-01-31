package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	TokenLedgerRepository interface {
		Create(ctx context.Context, tx *gorm.DB, ledger entity.TokenLedger) (entity.TokenLedger, error)

		// FOR USER
		GetWalletHistory(ctx context.Context, tx *gorm.DB, WalletID string, limit int, offset int) ([]entity.TokenLedger, error)
		GetByWalletPeriod(ctx context.Context, tx *gorm.DB, WalletID string, startDate string, endDate string) ([]entity.TokenLedger, error)

		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenLedger, error)
		CountByWalletID(ctx context.Context, tx *gorm.DB, WalletID string) (int, error)
	}

	tokenLedgerRepository struct {
		db *gorm.DB
	}
)

func NewTokenLedgerRepository(db *gorm.DB) TokenLedgerRepository {
	return &tokenLedgerRepository{db: db}
}

func (r *tokenLedgerRepository) Create(ctx context.Context, tx *gorm.DB, ledger entity.TokenLedger) (entity.TokenLedger, error) {
	if tx == nil {
		tx = r.db
	}
	err := tx.WithContext(ctx).Create(&ledger).Error
	if err != nil {
		return entity.TokenLedger{}, err
	}
	return ledger, nil
}

func (r *tokenLedgerRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenLedger, error) {
	if tx == nil {
		tx = r.db
	}
	var ledger entity.TokenLedger
	err := tx.WithContext(ctx).First(&ledger, id).Error
	if err != nil {
		return entity.TokenLedger{}, err
	}
	return ledger, nil
}

func (r *tokenLedgerRepository) GetWalletHistory(ctx context.Context, tx *gorm.DB, WalletID string, limit int, offset int) ([]entity.TokenLedger, error) {
	if tx == nil {
		tx = r.db
	}
	var ledgers []entity.TokenLedger
	err := tx.WithContext(ctx).Where("wallet_id = ?", WalletID).Limit(limit).Offset(offset).Find(&ledgers).Error
	if err != nil {
		return nil, err
	}
	return ledgers, nil
}

func (r *tokenLedgerRepository) GetByWalletPeriod(ctx context.Context, tx *gorm.DB, WalletID string, startDate string, endDate string) ([]entity.TokenLedger, error) {
	if tx == nil {
		tx = r.db
	}
	var ledgers []entity.TokenLedger
	err := tx.WithContext(ctx).Where("wallet_id = ? AND created_at BETWEEN ? AND ?", WalletID, startDate, endDate).Find(&ledgers).Error
	if err != nil {
		return nil, err
	}
	return ledgers, nil
}

func (r *tokenLedgerRepository) CountByWalletID(ctx context.Context, tx *gorm.DB, WalletID string) (int, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	err := tx.WithContext(ctx).Model(&entity.TokenLedger{}).Where("wallet_id = ?", WalletID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
