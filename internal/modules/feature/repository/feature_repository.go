package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	FeatureRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.Feature) (entity.Feature, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Feature, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Feature, error)
	}

	featureRepository struct {
		db *gorm.DB
	}
)

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepository{
		db: db,
	}
}

func (r *featureRepository) Create(ctx context.Context, tx *gorm.DB, model entity.Feature) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(model).Error; err != nil {
		return entity.Feature{}, err
	}

	return model, nil
}

func (r *featureRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.Feature
	if err := tx.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *featureRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}

	var model entity.Feature
	if err := tx.WithContext(ctx).Where(&model, "id = ?", id).Error; err != nil {
		return entity.Feature{}, err
	}

	return model, nil
}
