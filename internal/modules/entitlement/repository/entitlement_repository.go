package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	EntitlementRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.Entitlement) (entity.Entitlement, error)
		GetAll(ctx context.Context, tx *gorm.DB, featureID *uuid.UUID) ([]entity.Entitlement, error)
		GetAllWithMappingCounts(ctx context.Context, tx *gorm.DB, featureID *uuid.UUID) ([]entity.Entitlement, map[uuid.UUID]int64, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Entitlement, error)
		GetByKey(ctx context.Context, tx *gorm.DB, key string) (entity.Entitlement, error)
		GetByComponents(ctx context.Context, tx *gorm.DB, featureID uuid.UUID, subFeatureID *uuid.UUID, actionCategoryID uuid.UUID) (entity.Entitlement, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	entitlementRepository struct {
		db *gorm.DB
	}
)

func NewEntitlementRepository(db *gorm.DB) EntitlementRepository {
	return &entitlementRepository{db: db}
}

func (r *entitlementRepository) Create(ctx context.Context, tx *gorm.DB, model entity.Entitlement) (entity.Entitlement, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.Entitlement{}, err
	}
	return model, nil
}

// GetAll mengambil semua entitlement, opsional difilter berdasarkan feature_id.
// Selalu preload Feature, SubFeature, dan ActionCategory untuk kebutuhan response mapping.
func (r *entitlementRepository) GetAll(ctx context.Context, tx *gorm.DB, featureID *uuid.UUID) ([]entity.Entitlement, error) {
	if tx == nil {
		tx = r.db
	}
	query := tx.WithContext(ctx).
		Preload("Feature").
		Preload("SubFeature").
		Preload("ActionCategory")

	if featureID != nil {
		query = query.Where("feature_id = ?", *featureID)
	}

	var models []entity.Entitlement
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// GetAllWithMappingCounts mengambil semua entitlement beserta jumlah DurationAccessMapping
// yang menggunakannya, dalam dua query flat — menghindari N+1.
func (r *entitlementRepository) GetAllWithMappingCounts(ctx context.Context, tx *gorm.DB, featureID *uuid.UUID) ([]entity.Entitlement, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).
		Preload("Feature").
		Preload("SubFeature").
		Preload("ActionCategory")

	if featureID != nil {
		query = query.Where("feature_id = ?", *featureID)
	}

	var models []entity.Entitlement
	if err := query.Find(&models).Error; err != nil {
		return nil, nil, err
	}

	// Kumpulkan semua ID entitlement untuk query aggregate
	ids := make([]uuid.UUID, len(models))
	for i, m := range models {
		ids[i] = m.ID
	}

	type countResult struct {
		EntitlementID uuid.UUID `gorm:"column:entitlement_id"`
		Count         int64     `gorm:"column:count"`
	}
	var counts []countResult
	if len(ids) > 0 {
		if err := tx.WithContext(ctx).Model(&entity.DurationAccessMapping{}).
			Select("entitlement_id, COUNT(*) as count").
			Where("entitlement_id IN ?", ids).
			Group("entitlement_id").
			Scan(&counts).Error; err != nil {
			return nil, nil, err
		}
	}

	countMap := make(map[uuid.UUID]int64)
	for _, c := range counts {
		countMap[c.EntitlementID] = c.Count
	}

	return models, countMap, nil
}

func (r *entitlementRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Entitlement, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Entitlement
	if err := tx.WithContext(ctx).
		Preload("Feature").
		Preload("SubFeature").
		Preload("ActionCategory").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.Entitlement{}, err
	}
	return model, nil
}

func (r *entitlementRepository) GetByKey(ctx context.Context, tx *gorm.DB, key string) (entity.Entitlement, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Entitlement
	if err := tx.WithContext(ctx).
		Preload("Feature").
		Preload("SubFeature").
		Preload("ActionCategory").
		First(&model, "key = ?", key).Error; err != nil {
		return entity.Entitlement{}, err
	}
	return model, nil
}

// GetByComponents digunakan untuk cek duplikasi sebelum Create —
// kombinasi (feature_id, sub_feature_id, action_category_id) harus unik.
func (r *entitlementRepository) GetByComponents(ctx context.Context, tx *gorm.DB, featureID uuid.UUID, subFeatureID *uuid.UUID, actionCategoryID uuid.UUID) (entity.Entitlement, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.Entitlement
	query := tx.WithContext(ctx).Where("feature_id = ? AND action_category_id = ?", featureID, actionCategoryID)

	if subFeatureID != nil {
		query = query.Where("sub_feature_id = ?", *subFeatureID)
	} else {
		query = query.Where("sub_feature_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		return entity.Entitlement{}, err
	}
	return model, nil
}

func (r *entitlementRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.Entitlement{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CountMappings menghitung berapa DurationAccessMapping yang menggunakan entitlement ini.
// Dipakai oleh service sebagai guard sebelum Delete.
func (r *entitlementRepository) CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.DurationAccessMapping{}).
		Where("entitlement_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
