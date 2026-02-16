package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	HakAksesRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.HakAkses) (entity.HakAkses, error)
		FindByFeatureAction(ctx context.Context, tx *gorm.DB, feature, action string) (entity.HakAkses, error)
	}

	hakAksesRepository struct {
		db *gorm.DB
	}
)

func NewHakAksesRepository(db *gorm.DB) HakAksesRepository {
	return &hakAksesRepository{
		db: db,
	}
}

func (r *hakAksesRepository) Create(ctx context.Context, tx *gorm.DB, model entity.HakAkses) (entity.HakAkses, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(model).Error; err != nil {
		return entity.HakAkses{}, err
	}

	return model, nil
}

func (r *hakAksesRepository) FindByFeatureAction(ctx context.Context, tx *gorm.DB, feature, action string) (entity.HakAkses, error) {
	if tx == nil {
		tx = r.db
	}

	var model entity.HakAkses
	if err := tx.WithContext(ctx).Preload("MembershipPlan").First(&model, "feature = ? AND action = ?", feature, action).Error; err != nil {
		return entity.HakAkses{}, err
	}

	return model, nil
}
