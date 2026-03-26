package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PromoCodeUsageRepository interface {
		Create(ctx context.Context, tx *gorm.DB, promoCodeUsage entity.PromoCodeUsage) (entity.PromoCodeUsage, error)
		GetAllPromoCodeUsage(ctx context.Context, tx *gorm.DB) ([]entity.PromoCodeUsage, error)
		IsPromoCodeUsedByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (bool, error)
	}

	promoCodeUsageRepository struct {
		db *gorm.DB
	}
)

func NewPromoCodeUsageRepository(db *gorm.DB) PromoCodeUsageRepository {
	return &promoCodeUsageRepository{
		db: db,
	}
}

func (r *promoCodeUsageRepository) Create(ctx context.Context, tx *gorm.DB, promoCodeUsage entity.PromoCodeUsage) (entity.PromoCodeUsage, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&promoCodeUsage).Error; err != nil {
		return entity.PromoCodeUsage{}, err
	}

	return promoCodeUsage, nil
}

func (r *promoCodeUsageRepository) GetAllPromoCodeUsage(ctx context.Context, tx *gorm.DB) ([]entity.PromoCodeUsage, error) {
	if tx == nil {
		tx = r.db
	}

	var promoCodeUsages []entity.PromoCodeUsage
	if err := tx.WithContext(ctx).Find(&promoCodeUsages).Error; err != nil {
		return nil, err
	}

	return promoCodeUsages, nil
}

func (r *promoCodeUsageRepository) IsPromoCodeUsedByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (bool, error) {
	if tx == nil {
		tx = r.db
	}

	var count int64

	if err := tx.WithContext(ctx).Model(&entity.PromoCodeUsage{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
