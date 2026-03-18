package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	GetTransactionsFilter struct {
		Page     int
		PageSize int
		Status   string
	}

	PaymentTransactionsRepository interface {
		Create(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error)
		GetByID(ctx context.Context, tx *gorm.DB, paymentId uuid.UUID) (entity.PaymentTransactions, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.PaymentTransactions, error)
		GetByPaymentInvoiceID(ctx context.Context, tx *gorm.DB, paymentInvoiceId string) (entity.PaymentTransactions, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.PaymentTransactions, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.PaymentTransactions, error)
		Update(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error)

		// Checkout specific methods
		GetTransactionByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error)
		GetTransactionsByUserIDWithFilter(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filter GetTransactionsFilter) ([]entity.PaymentTransactions, int64, error)
		GetPendingTransactionByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.PaymentTransactions, error)
		GetMembershipByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error)
		GetPlanDurationByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipDuration, error)
		GetTokenBundleByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.TokenBundlePackage, error)
		GetActiveTokenBundles(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error)
	}

	paymentTransactionRepository struct {
		db *gorm.DB
	}
)

func NewPaymentTransactionRepository(db *gorm.DB) PaymentTransactionsRepository {
	return &paymentTransactionRepository{
		db: db,
	}
}

func (r *paymentTransactionRepository) Create(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&paymentTransaction).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}

	return paymentTransaction, nil
}

func (r *paymentTransactionRepository) GetByID(ctx context.Context, tx *gorm.DB, paymentId uuid.UUID) (entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	var paymentTransaction entity.PaymentTransactions
	if err := tx.WithContext(ctx).First(&paymentTransaction, "id = ?", paymentId).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}

	return paymentTransaction, nil
}

func (r *paymentTransactionRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	var userPayment []entity.PaymentTransactions
	if err := tx.WithContext(ctx).Find(&userPayment, "user_id = ?", userId).Error; err != nil {
		return []entity.PaymentTransactions{}, err
	}

	return userPayment, nil
}

func (r *paymentTransactionRepository) GetByPaymentInvoiceID(ctx context.Context, tx *gorm.DB, paymentInvoiceId string) (entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	var payment entity.PaymentTransactions
	if err := tx.WithContext(ctx).First(&payment, "payment_external_id = ?", paymentInvoiceId).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}

	return payment, nil
}

func (r *paymentTransactionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	var payments []entity.PaymentTransactions
	if err := tx.WithContext(ctx).Find(&payments).Error; err != nil {
		return []entity.PaymentTransactions{}, err
	}

	return payments, nil
}

func (r *paymentTransactionRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	offset := (page + 1) * totalPage

	var payments []entity.PaymentTransactions
	if err := tx.WithContext(ctx).Find(&payments).Offset(offset).Limit(totalPage).Error; err != nil {
		return []entity.PaymentTransactions{}, err
	}

	return payments, nil
}

func (r *paymentTransactionRepository) Update(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&paymentTransaction).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}

	return paymentTransaction, nil
}

// ─── Checkout Specific Methods ────────────────────────────────────────────────

func (r *paymentTransactionRepository) GetTransactionByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.PaymentTransactions, error) {
	if tx == nil {tx = r.db}
	var model entity.PaymentTransactions
	if err := tx.WithContext(ctx).
		Preload("Plan").
		Preload("Duration").
		Preload("TokenBundlePackage").
		First(&model, "transaction_id = ?", transactionID).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}
	return model, nil
}

func (r *paymentTransactionRepository) GetTransactionsByUserIDWithFilter(ctx context.Context, tx *gorm.DB, userID uuid.UUID, filter GetTransactionsFilter) ([]entity.PaymentTransactions, int64, error) {
	if tx == nil {tx = r.db}

	query := tx.WithContext(ctx).Model(&entity.PaymentTransactions{}).Where("user_id = ?", userID)

	if filter.Status != "" {
		query = query.Where("payment_status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := filter.Page * filter.PageSize
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	var models []entity.PaymentTransactions
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

func (r *paymentTransactionRepository) GetPendingTransactionByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.PaymentTransactions, error) {
	if tx == nil {tx = r.db}
	var model entity.PaymentTransactions
	if err := tx.WithContext(ctx).
		Where("user_id = ? AND payment_status = ?", userID, entity.PaymentStatusPending).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		return entity.PaymentTransactions{}, err
	}
	return model, nil
}

func (r *paymentTransactionRepository) GetMembershipByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error) {
	if tx == nil {tx = r.db}
	var model entity.Memberships
	if err := tx.WithContext(ctx).
		Preload("Duration").
		First(&model, "user_id = ?", userID).Error; err != nil {
		return entity.Memberships{}, err
	}
	return model, nil
}

func (r *paymentTransactionRepository) GetPlanDurationByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipDuration, error) {
	if tx == nil {tx = r.db}
	var model entity.MembershipDuration
	// In the real system, MembershipDuration has no explicit Plan relation column natively.
	// We'll just fetch by ID. 
	if err := tx.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return entity.MembershipDuration{}, err
	}
	return model, nil
}

func (r *paymentTransactionRepository) GetTokenBundleByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.TokenBundlePackage, error) {
	if tx == nil {tx = r.db}
	var model entity.TokenBundlePackage
	if err := tx.WithContext(ctx).First(&model, "id = ? AND is_active = true", id).Error; err != nil {
		return entity.TokenBundlePackage{}, err
	}
	return model, nil
}

func (r *paymentTransactionRepository) GetActiveTokenBundles(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error) {
	if tx == nil {tx = r.db}
	var models []entity.TokenBundlePackage
	if err := tx.WithContext(ctx).
		Where("is_active = true").
		Order("display_order ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}
