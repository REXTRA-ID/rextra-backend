package repository

import (
	"context"
	"rextra-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopupTransactionFilters struct {
	Status    *string
	Type      *string
	StartDate *string // Expects "YYYY-MM-DD"
	EndDate   *string // Expects "YYYY-MM-DD"
}

type (
	TopupTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, data *entity.TopupTransaction) (*entity.TopupTransaction, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.TopupTransaction, error)
		GetByInvoiceID(ctx context.Context, tx *gorm.DB, invoiceID string) (*entity.TopupTransaction, error)
		GetUserTransactions(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filters TopupTransactionFilters, limit, offset int) ([]entity.TopupTransaction, int64, error)
		UpdateStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status entity.TopupStatus, paidAt *time.Time, metadata map[string]any) error
		SetLedgerID(ctx context.Context, tx *gorm.DB, id uuid.UUID, ledgerID uuid.UUID) error
	}

	topupTransactionRepository struct {
		db *gorm.DB
	}
)

func NewTopupTransactionRepository(db *gorm.DB) TopupTransactionRepository {
	return &topupTransactionRepository{
		db: db,
	}
}

func (r *topupTransactionRepository) Create(ctx context.Context, tx *gorm.DB, data *entity.TopupTransaction) (*entity.TopupTransaction, error) {
	if tx == nil {
		tx = r.db
	}
	db := tx.WithContext(ctx)
	if err := db.Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *topupTransactionRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.TopupTransaction, error) {
	if tx == nil {
		tx = r.db
	}
	var transaction entity.TopupTransaction
	db := tx.WithContext(ctx)
	if err := db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *topupTransactionRepository) GetByInvoiceID(ctx context.Context, tx *gorm.DB, invoiceID string) (*entity.TopupTransaction, error) {
	if tx == nil {
		tx = r.db
	}
	var transaction entity.TopupTransaction
	db := tx.WithContext(ctx)
	if err := db.Where("invoice_id = ?", invoiceID).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *topupTransactionRepository) GetUserTransactions(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filters TopupTransactionFilters, limit, offset int) ([]entity.TopupTransaction, int64, error) {
	if tx == nil {
		tx = r.db
	}
	var transactions []entity.TopupTransaction
	var total int64

	query := tx.WithContext(ctx).Model(&entity.TopupTransaction{}).Where("user_id = ?", userID)
	query = r.applyTransactionFilters(query, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("BundlePackage").Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *topupTransactionRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status entity.TopupStatus, paidAt *time.Time, metadata map[string]any) error {
	if tx == nil {
		tx = r.db
	}
	db := tx.WithContext(ctx)
	updates := map[string]interface{}{
		"status":  status,
		"paid_at": paidAt,
	}
	if metadata != nil {
		updates["metadata"] = metadata
	}
	return db.Model(&entity.TopupTransaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *topupTransactionRepository) SetLedgerID(ctx context.Context, tx *gorm.DB, id uuid.UUID, ledgerID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	db := tx.WithContext(ctx)
	return db.Model(&entity.TopupTransaction{}).Where("id = ?", id).Update("ledger_id", ledgerID).Error
}

// Helper
func (r *topupTransactionRepository) applyTransactionFilters(query *gorm.DB, filters TopupTransactionFilters) *gorm.DB {
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.StartDate != nil {
		if startDate, err := time.Parse("2006-01-02", *filters.StartDate); err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if filters.EndDate != nil {
		if endDate, err := time.Parse("2006-01-02", *filters.EndDate); err == nil {
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query = query.Where("created_at <= ?", endDate)
		}
	}
	return query
}
