package repository

import (
	"context"
	"rextra-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionsRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PaymentTransactions, error)
		GetByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error)
		GetByMerchantRef(ctx context.Context, tx *gorm.DB, merchantRef string) (entity.PaymentTransactions, error)
		GetAllPaginatedWithFilter(ctx context.Context, tx *gorm.DB, offset, limit int, userID uuid.UUID, status, planName, search string, dateFrom, dateTo time.Time) ([]entity.PaymentTransactions, int64, error)
		GetByUserIDPaginated(ctx context.Context, tx *gorm.DB, userID uuid.UUID, offset, limit int, status string) ([]entity.PaymentTransactions, int64, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error)
	}

	paymentTransactionsRepository struct {
		db *gorm.DB
	}
)

func NewPaymentTransactionRepository(db *gorm.DB) PaymentTransactionsRepository {
	return &paymentTransactionsRepository{db: db}
}

func (r *paymentTransactionsRepository) Create(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *paymentTransactionsRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var result entity.PaymentTransactions
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").First(&result, "id = ?", id).Error; err != nil { return entity.PaymentTransactions{}, err }
	return result, nil
}

func (r *paymentTransactionsRepository) GetByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var result entity.PaymentTransactions
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").First(&result, "transaction_id = ?", transactionID).Error; err != nil { return entity.PaymentTransactions{}, err }
	return result, nil
}

func (r *paymentTransactionsRepository) GetByMerchantRef(ctx context.Context, tx *gorm.DB, merchantRef string) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var result entity.PaymentTransactions
	if err := tx.WithContext(ctx).First(&result, "merchant_ref = ?", merchantRef).Error; err != nil { return entity.PaymentTransactions{}, err }
	return result, nil
}

func (r *paymentTransactionsRepository) GetAllPaginatedWithFilter(ctx context.Context, tx *gorm.DB, offset, limit int, userID uuid.UUID, status, planName, search string, dateFrom, dateTo time.Time) ([]entity.PaymentTransactions, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.PaymentTransactions{})
	if userID != uuid.Nil { query = query.Where("payment_transactions.user_id = ?", userID) }
	if status != "" { query = query.Where("payment_transactions.payment_status = ?", status) }
	if planName != "" { query = query.Where("payment_transactions.to_plan = ?", planName) }
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("payment_transactions.user_name ILIKE ? OR payment_transactions.user_email ILIKE ?", pattern, pattern)
	}
	if !dateFrom.IsZero() { query = query.Where("payment_transactions.created_at >= ?", dateFrom) }
	if !dateTo.IsZero() { query = query.Where("payment_transactions.created_at < ?", dateTo.AddDate(0, 0, 1)) }
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.PaymentTransactions
	if err := query.Select("payment_transactions.*").Order("payment_transactions.created_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *paymentTransactionsRepository) GetByUserIDPaginated(ctx context.Context, tx *gorm.DB, userID uuid.UUID, offset, limit int, status string) ([]entity.PaymentTransactions, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.PaymentTransactions{}).Where("user_id = ?", userID)
	if status != "" { query = query.Where("payment_status = ?", status) }
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.PaymentTransactions
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *paymentTransactionsRepository) Update(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}
