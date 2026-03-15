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
		GetAllWithCounts(ctx context.Context, tx *gorm.DB) ([]entity.Feature, map[uuid.UUID]int64, map[uuid.UUID]int64, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Feature, error)
		GetBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Feature, error)
		GetByPrefix(ctx context.Context, tx *gorm.DB, prefix string) (entity.Feature, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.Feature) (entity.Feature, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
		CountSubFeatures(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	featureRepository struct {
		db *gorm.DB
	}
)

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
	return &featureRepository{db: db}
}

func (r *featureRepository) Create(ctx context.Context, tx *gorm.DB, model entity.Feature) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
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

func (r *featureRepository) GetAllWithCounts(ctx context.Context, tx *gorm.DB) ([]entity.Feature, map[uuid.UUID]int64, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.Feature
	if err := tx.WithContext(ctx).Preload("SubFeatures").Find(&models).Error; err != nil {
		return nil, nil, nil, err
	}

	type countResult struct {
		ID    uuid.UUID `gorm:"column:id"`
		Count int64     `gorm:"column:count"`
	}

	var featureCounts []countResult
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Select("feature_id as id, COUNT(*) as count").
		Group("feature_id").
		Scan(&featureCounts).Error; err != nil {
		return nil, nil, nil, err
	}
	featureCountMap := make(map[uuid.UUID]int64)
	for _, c := range featureCounts {
		featureCountMap[c.ID] = c.Count
	}

	var subFeatureCounts []countResult
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Select("sub_feature_id as id, COUNT(*) as count").
		Where("sub_feature_id IS NOT NULL").
		Group("sub_feature_id").
		Scan(&subFeatureCounts).Error; err != nil {
		return nil, nil, nil, err
	}
	subFeatureCountMap := make(map[uuid.UUID]int64)
	for _, c := range subFeatureCounts {
		subFeatureCountMap[c.ID] = c.Count
	}

	return models, featureCountMap, subFeatureCountMap, nil
}

func (r *featureRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Feature
	if err := tx.WithContext(ctx).
		Preload("SubFeatures").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.Feature{}, err
	}
	return model, nil
}

func (r *featureRepository) GetBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Feature
	if err := tx.WithContext(ctx).First(&model, "slug = ?", slug).Error; err != nil {
		return entity.Feature{}, err
	}
	return model, nil
}

func (r *featureRepository) GetByPrefix(ctx context.Context, tx *gorm.DB, prefix string) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Feature
	if err := tx.WithContext(ctx).First(&model, "prefix = ?", prefix).Error; err != nil {
		return entity.Feature{}, err
	}
	return model, nil
}

func (r *featureRepository) Update(ctx context.Context, tx *gorm.DB, model entity.Feature) (entity.Feature, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.Feature{}, err
	}
	return model, nil
}

func (r *featureRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.Feature{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *featureRepository) CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Where("feature_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *featureRepository) CountSubFeatures(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.SubFeature{}).
		Where("feature_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
