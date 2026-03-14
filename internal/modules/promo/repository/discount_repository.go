package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DiscountRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.Discounts) (entity.Discounts, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, status, appliesTo, search string) ([]entity.Discounts, int64, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Discounts, error)
		GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.Discounts, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.Discounts) (entity.Discounts, error)
		SoftDelete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		IncrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		DecrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	}

	discountRepository struct {
		db *gorm.DB
	}
)

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{db: db}
}

func (r *discountRepository) Create(ctx context.Context, tx *gorm.DB, model entity.Discounts) (entity.Discounts, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.Discounts{}, err }
	return model, nil
}

func (r *discountRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, status, appliesTo, search string) ([]entity.Discounts, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.Discounts{})
	if status != "" { query = query.Where("status = ?", status) }
	if appliesTo != "" { query = query.Where("applies_to = ?", appliesTo) }
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ?", pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.Discounts
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *discountRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Discounts, error) {
	if tx == nil { tx = r.db }
	var result entity.Discounts
	if err := tx.WithContext(ctx).First(&result, "id = ?", id).Error; err != nil { return entity.Discounts{}, err }
	return result, nil
}

func (r *discountRepository) GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.Discounts, error) {
	if tx == nil { tx = r.db }
	var result entity.Discounts
	if err := tx.WithContext(ctx).Where("UPPER(code) = UPPER(?)", code).First(&result).Error; err != nil { return entity.Discounts{}, err }
	return result, nil
}

func (r *discountRepository) Update(ctx context.Context, tx *gorm.DB, model entity.Discounts) (entity.Discounts, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil { return entity.Discounts{}, err }
	return model, nil
}

func (r *discountRepository) SoftDelete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Model(&entity.Discounts{}).Where("id = ?", id).Update("status", entity.DiscountStatusInactive).Error
}

func (r *discountRepository) IncrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Model(&entity.Discounts{}).Where("id = ?", id).UpdateColumn("current_redemptions", gorm.Expr("current_redemptions + 1")).Error
}

func (r *discountRepository) DecrementRedemption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Model(&entity.Discounts{}).Where("id = ?", id).UpdateColumn("current_redemptions", gorm.Expr("GREATEST(0, current_redemptions - 1)")).Error
}
