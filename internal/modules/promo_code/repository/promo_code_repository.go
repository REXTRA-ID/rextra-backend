package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	PromoCodeRepository interface {
		GetAllPromoCodes(ctx context.Context, tx *gorm.DB) ([]entity.PromoCodes, error)
		GetPromoCodeByCodeName(ctx context.Context, tx *gorm.DB, promoCode string) (*entity.PromoCodes, error)
	}

	promoCodeRepository struct {
		db *gorm.DB
	}
)

func NewPromoCodesRepository(db *gorm.DB) PromoCodeRepository {
	return &promoCodeRepository{
		db: db,
	}
}

func (r *promoCodeRepository) GetAllPromoCodes(ctx context.Context, tx *gorm.DB) ([]entity.PromoCodes, error) {
	if tx == nil {
		tx = r.db
	}

	var promoCodes []entity.PromoCodes
	if err := tx.WithContext(ctx).Find(&promoCodes).Error; err != nil {
		return nil, err
	}

	return promoCodes, nil
}

func (r *promoCodeRepository) GetPromoCodeByCodeName(ctx context.Context, tx *gorm.DB, promoCode string) (*entity.PromoCodes, error) {
	if tx == nil {
		tx = r.db
	}

	var promocode *entity.PromoCodes
	if err := tx.WithContext(ctx).First(promocode, "code = ?", promoCode).Error; err != nil {
		return nil, err
	}

	return promocode, nil
}
