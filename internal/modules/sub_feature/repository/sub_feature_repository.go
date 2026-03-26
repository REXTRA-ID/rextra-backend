package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	SubFeatureRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.SubFeature) (entity.SubFeature, error)
		GetAllByFeatureID(ctx context.Context, tx *gorm.DB, featureID uuid.UUID) ([]entity.SubFeature, error)
		GetAllByFeatureIDWithCounts(ctx context.Context, tx *gorm.DB, featureID uuid.UUID) ([]entity.SubFeature, map[uuid.UUID]int64, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.SubFeature, error)
		GetBySlugAndFeatureID(ctx context.Context, tx *gorm.DB, slug string, featureID uuid.UUID) (entity.SubFeature, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.SubFeature) (entity.SubFeature, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	subFeatureRepository struct {
		db *gorm.DB
	}
)

func NewSubFeatureRepository(db *gorm.DB) SubFeatureRepository {
	return &subFeatureRepository{db: db}
}

func (r *subFeatureRepository) Create(ctx context.Context, tx *gorm.DB, model entity.SubFeature) (entity.SubFeature, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.SubFeature{}, err
	}
	return model, nil
}

func (r *subFeatureRepository) GetAllByFeatureID(ctx context.Context, tx *gorm.DB, featureID uuid.UUID) ([]entity.SubFeature, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.SubFeature
	if err := tx.WithContext(ctx).
		Where("feature_id = ?", featureID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// GetAllByFeatureIDWithCounts mengambil semua SubFeature milik satu Feature
// beserta entitlement count masing-masing dalam dua query flat.
func (r *subFeatureRepository) GetAllByFeatureIDWithCounts(ctx context.Context, tx *gorm.DB, featureID uuid.UUID) ([]entity.SubFeature, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.SubFeature
	if err := tx.WithContext(ctx).
		Where("feature_id = ?", featureID).
		Find(&models).Error; err != nil {
		return nil, nil, err
	}

	type countResult struct {
		SubFeatureID uuid.UUID `gorm:"column:sub_feature_id"`
		Count        int64     `gorm:"column:count"`
	}
	var counts []countResult
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Select("sub_feature_id, COUNT(*) as count").
		Where("sub_feature_id IS NOT NULL AND feature_id = ?", featureID).
		Group("sub_feature_id").
		Scan(&counts).Error; err != nil {
		return nil, nil, err
	}

	countMap := make(map[uuid.UUID]int64)
	for _, c := range counts {
		countMap[c.SubFeatureID] = c.Count
	}

	return models, countMap, nil
}

func (r *subFeatureRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.SubFeature, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.SubFeature
	if err := tx.WithContext(ctx).
		Preload("Feature").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.SubFeature{}, err
	}
	return model, nil
}

func (r *subFeatureRepository) GetBySlugAndFeatureID(ctx context.Context, tx *gorm.DB, slug string, featureID uuid.UUID) (entity.SubFeature, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.SubFeature
	if err := tx.WithContext(ctx).
		First(&model, "slug = ? AND feature_id = ?", slug, featureID).Error; err != nil {
		return entity.SubFeature{}, err
	}
	return model, nil
}

func (r *subFeatureRepository) Update(ctx context.Context, tx *gorm.DB, model entity.SubFeature) (entity.SubFeature, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.SubFeature{}, err
	}
	return model, nil
}

func (r *subFeatureRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.SubFeature{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *subFeatureRepository) CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Where("sub_feature_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
