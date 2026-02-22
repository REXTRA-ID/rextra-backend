package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	ActionCategoryRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.ActionCategory, error)
	}

	actionCategoryRepository struct {
		db *gorm.DB
	}
)

func NewActionCategoryRepository(db *gorm.DB) ActionCategoryRepository {
	return &actionCategoryRepository{
		db: db,
	}
}

func (r *actionCategoryRepository) Create(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(model).Error; err != nil {
		return entity.ActionCategory{}, err
	}

	return model, nil
}

func (r *actionCategoryRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.ActionCategory
	if err := tx.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}
