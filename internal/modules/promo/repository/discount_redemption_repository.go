package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DiscountRedemptionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.DiscountRedemption) (entity.DiscountRedemption, error)
		GetByDiscountIDPaginated(ctx context.Context, tx *gorm.DB, discountID uuid.UUID, offset, limit int, status string) ([]entity.DiscountRedemption, int64, error)
		CountByUserAndDiscount(ctx context.Context, tx *gorm.DB, userID, discountID uuid.UUID) (int64, error)
		GetByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.DiscountRedemption, error)
		Reverse(ctx context.Context, tx *gorm.DB, transactionID, reason string) error
		Update(ctx context.Context, tx *gorm.DB, model entity.DiscountRedemption) (entity.DiscountRedemption, error)
	}

	discountRedemptionRepository struct {
		db *gorm.DB
	}
)

func NewDiscountRedemptionRepository(db *gorm.DB) DiscountRedemptionRepository {
	return &discountRedemptionRepository{db: db}
}

func (r *discountRedemptionRepository) Create(ctx context.Context, tx *gorm.DB, model entity.DiscountRedemption) (entity.DiscountRedemption, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.DiscountRedemption{}, err }
	return model, nil
}

func (r *discountRedemptionRepository) GetByDiscountIDPaginated(ctx context.Context, tx *gorm.DB, discountID uuid.UUID, offset, limit int, status string) ([]entity.DiscountRedemption, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.DiscountRedemption{}).Where("discount_id = ?", discountID)
	if status != "" { query = query.Where("status = ?", status) }
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.DiscountRedemption
	if err := query.Order("applied_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *discountRedemptionRepository) CountByUserAndDiscount(ctx context.Context, tx *gorm.DB, userID, discountID uuid.UUID) (int64, error) {
	if tx == nil { tx = r.db }
	var count int64
	err := tx.WithContext(ctx).Model(&entity.DiscountRedemption{}).Where("user_id = ? AND discount_id = ? AND status = ?", userID, discountID, entity.RedemptionStatusApplied).Count(&count).Error
	return count, err
}

func (r *discountRedemptionRepository) GetByTransactionID(ctx context.Context, tx *gorm.DB, transactionID string) (entity.DiscountRedemption, error) {
	if tx == nil { tx = r.db }
	var result entity.DiscountRedemption
	if err := tx.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&result).Error; err != nil { return entity.DiscountRedemption{}, err }
	return result, nil
}

func (r *discountRedemptionRepository) Reverse(ctx context.Context, tx *gorm.DB, transactionID, reason string) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Model(&entity.DiscountRedemption{}).Where("transaction_id = ? AND status = ?", transactionID, entity.RedemptionStatusApplied).Updates(map[string]interface{}{
		"status":         entity.RedemptionStatusReversed,
		"reversed_at":    gorm.Expr("NOW()"),
		"reverse_reason": reason,
	}).Error
}

func (r *discountRedemptionRepository) Update(ctx context.Context, tx *gorm.DB, model entity.DiscountRedemption) (entity.DiscountRedemption, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil { return entity.DiscountRedemption{}, err }
	return model, nil
}
