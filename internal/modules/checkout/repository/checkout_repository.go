package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CheckoutRepository interface {
		CreateTransaction(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error)
		GetTransactionByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PaymentTransactions, error)
		GetTransactionByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error)
		GetTransactionsByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filter GetTransactionsFilter) ([]entity.PaymentTransactions, int64, error)
		GetPendingTransactionByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.PaymentTransactions, error)
		UpdateTransaction(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error)
		GetMembershipByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error)
		GetPlanDurationByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PlanDuration, error)
		GetTokenBundleByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.TokenBundlePackage, error)
		GetActiveTokenBundles(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error)
	}

	GetTransactionsFilter struct {
		Page     int
		PageSize int
		Status   string
	}

	checkoutRepository struct {
		db *gorm.DB
	}
)

func NewCheckoutRepository(db *gorm.DB) CheckoutRepository {
	return &checkoutRepository{db: db}
}

func (r *checkoutRepository) CreateTransaction(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *checkoutRepository) GetTransactionByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var model entity.PaymentTransactions
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").Preload("TokenBundlePackage").First(&model, "id = ?", id).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *checkoutRepository) GetTransactionByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var model entity.PaymentTransactions
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").Preload("TokenBundlePackage").First(&model, "transaction_id = ?", transactionID).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *checkoutRepository) GetTransactionsByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filter GetTransactionsFilter) ([]entity.PaymentTransactions, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.PaymentTransactions{}).Where("user_id = ?", userID)
	if filter.Status != "" { query = query.Where("payment_status = ?", filter.Status) }
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	offset := filter.Page * filter.PageSize
	if filter.PageSize == 0 { filter.PageSize = 20 }
	var models []entity.PaymentTransactions
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&models).Error; err != nil { return nil, 0, err }
	return models, total, nil
}

func (r *checkoutRepository) GetPendingTransactionByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	var model entity.PaymentTransactions
	if err := tx.WithContext(ctx).Where("user_id = ? AND payment_status = ?", userID, entity.PaymentStatusPending).Order("created_at DESC").First(&model).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *checkoutRepository) UpdateTransaction(ctx context.Context, tx *gorm.DB, model entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil { return entity.PaymentTransactions{}, err }
	return model, nil
}

func (r *checkoutRepository) GetMembershipByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error) {
	if tx == nil { tx = r.db }
	var model entity.Memberships
	if err := tx.WithContext(ctx).Preload("Duration").First(&model, "user_id = ?", userID).Error; err != nil { return entity.Memberships{}, err }
	return model, nil
}

func (r *checkoutRepository) GetPlanDurationByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PlanDuration, error) {
	if tx == nil { tx = r.db }
	var model entity.PlanDuration
	if err := tx.WithContext(ctx).Preload("Plan").First(&model, "id = ?", id).Error; err != nil { return entity.PlanDuration{}, err }
	return model, nil
}

func (r *checkoutRepository) GetTokenBundleByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.TokenBundlePackage, error) {
	if tx == nil { tx = r.db }
	var model entity.TokenBundlePackage
	if err := tx.WithContext(ctx).First(&model, "id = ? AND is_active = true", id).Error; err != nil { return entity.TokenBundlePackage{}, err }
	return model, nil
}

func (r *checkoutRepository) GetActiveTokenBundles(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error) {
	if tx == nil { tx = r.db }
	var models []entity.TokenBundlePackage
	if err := tx.WithContext(ctx).Where("is_active = true").Order("display_order ASC").Find(&models).Error; err != nil { return nil, err }
	return models, nil
}
