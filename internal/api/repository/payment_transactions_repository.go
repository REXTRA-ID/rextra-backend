package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PaymentTransactionsRepository interface {
		Create(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error)
		GetByID(ctx context.Context, tx *gorm.DB, paymentId uuid.UUID) (entity.PaymentTransactions, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) ([]entity.PaymentTransactions, error)
		GetByXenditID(ctx context.Context, tx *gorm.DB, xenditId string) (entity.PaymentTransactions, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.PaymentTransactions, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.PaymentTransactions, error)
		Update(ctx context.Context, tx *gorm.DB, paymentTransaction entity.PaymentTransactions) (entity.PaymentTransactions, error)
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

func (r *paymentTransactionRepository) GetByXenditID(ctx context.Context, tx *gorm.DB, xenditId string) (entity.PaymentTransactions, error) {
	if tx == nil {
		tx = r.db
	}

	var payment entity.PaymentTransactions
	if err := tx.WithContext(ctx).First(&payment, "xendit_invoice_id = ?", xenditId).Error; err != nil {
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
