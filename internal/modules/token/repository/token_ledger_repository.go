package repository

import (
	"context"
	"rextra-backend/internal/entity"
	"time"

	"gorm.io/gorm"
)

type TokenLedgerFilter struct {
	Direction  *string
	SourceType *string
	StartDate  *string
	EndDate    *string
}

type (
	TokenLedgerRepository interface {
		Create(ctx context.Context, tx *gorm.DB, ledger entity.TokenLedger) (entity.TokenLedger, error)

		// FOR USER
		GetWalletHistory(ctx context.Context, tx *gorm.DB, WalletID string, limit int, offset int) ([]entity.TokenLedger, error)
		GetByWalletPeriod(ctx context.Context, tx *gorm.DB, WalletID string, startDate string, endDate string) ([]entity.TokenLedger, error)
		GetDetail(ctx context.Context, tx *gorm.DB, WalletID string, id string) (entity.TokenLedger, error)

		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenLedger, error)
		CountByWalletID(ctx context.Context, tx *gorm.DB, WalletID string) (int, error)

		GetAllActivity(ctx context.Context, tx *gorm.DB, filter TokenLedgerFilter, limit int, offset int, preloads ...string) ([]entity.TokenLedger, int64, error)
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

func (r *tokenLedgerRepository) GetDetail(ctx context.Context, tx *gorm.DB, WalletID string, id string) (entity.TokenLedger, error) {
	if tx == nil {
		tx = r.db
	}
	var ledger entity.TokenLedger
	err := tx.WithContext(ctx).Where("id = ? AND wallet_id = ?", id, WalletID).First(&ledger).Error
	if err != nil {
		return entity.TokenLedger{}, err
	}
	return ledger, nil
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

func (r *tokenLedgerRepository) GetAllActivity(ctx context.Context, tx *gorm.DB, filter TokenLedgerFilter, limit int, offset int, preloads ...string) ([]entity.TokenLedger, int64, error) {
	if tx == nil {
		tx = r.db
	}
	var ledgers []entity.TokenLedger
	var total int64
	query := tx.WithContext(ctx).Model(&entity.TokenLedger{})

	if filter.Direction != nil {
		query = query.Where("direction = ?", *filter.Direction)
	}
	if filter.SourceType != nil {
		query = query.Where("source_type = ?", *filter.SourceType)
	}
	if filter.StartDate != nil {
		if startDate, err := time.Parse("2006-01-02", *filter.StartDate); err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if filter.EndDate != nil {
		if endDate, err := time.Parse("2006-01-02", *filter.EndDate); err == nil {
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query = query.Where("created_at <= ?", endDate)
		}
	}

	// Count the total records that match the filter
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloads for the final data retrieval
	for _, p := range preloads {
		query = query.Preload(p)
	}

	// Apply order and pagination and find the records
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ledgers).Error
	if err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}
