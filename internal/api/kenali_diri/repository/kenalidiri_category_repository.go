package repository

import (
	"context"
	"encoding/json"
	"time"

	"rextra-backend/internal/entity"
	"rextra-backend/internal/pkg/cache"

	"gorm.io/gorm"
)

const (
	categoryCacheKeyAll = "kenalidiri_categories:all"
)

type (
	KenalidiriCategoryRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.KenaliDiriCategory, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.KenaliDiriCategory, error)
		GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.KenaliDiriCategory, error)
	}

	kenalidiriCategoryRepository struct {
		db    *gorm.DB
		cache cache.CacheService
		ttl   time.Duration
	}
)

func NewKenalidiriCategory(db *gorm.DB, cacheSvc cache.CacheService) KenalidiriCategoryRepository {
	return &kenalidiriCategoryRepository{
		db:    db,
		cache: cacheSvc,
		ttl:   24 * time.Hour,
	}
}

func (r *kenalidiriCategoryRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.KenaliDiriCategory, error) {
	if tx == nil {
		tx = r.db
	}

	if r.cache != nil {
		if cached, err := r.cache.Get(ctx, categoryCacheKeyAll); err == nil {
			if categories, ok := decodeCategoryCache(cached); ok {
				return categories, nil
			}
		}
	}

	var categories []entity.KenaliDiriCategory
	if err := tx.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&categories).Error; err != nil {
		return nil, err
	}

	if r.cache != nil {
		if bytes, err := json.Marshal(categories); err == nil {
			_ = r.cache.Set(ctx, categoryCacheKeyAll, bytes, r.ttl)
		}
	}

	return categories, nil
}

func (r *kenalidiriCategoryRepository) GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.KenaliDiriCategory, error) {
	if tx == nil {
		tx = r.db
	}

	var category entity.KenaliDiriCategory
	if err := tx.WithContext(ctx).Take(&category, "id = ?", id).Error; err != nil {
		return entity.KenaliDiriCategory{}, err
	}

	return category, nil
}

func (r *kenalidiriCategoryRepository) GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.KenaliDiriCategory, error) {
	if tx == nil {
		tx = r.db
	}

	var category entity.KenaliDiriCategory
	if err := tx.WithContext(ctx).Take(&category, "category_code = ?", code).Error; err != nil {
		return entity.KenaliDiriCategory{}, err
	}

	return category, nil
}

func decodeCategoryCache(data interface{}) ([]entity.KenaliDiriCategory, bool) {
	switch v := data.(type) {
	case []entity.KenaliDiriCategory:
		return v, true
	case []byte:
		var categories []entity.KenaliDiriCategory
		if err := json.Unmarshal(v, &categories); err == nil {
			return categories, true
		}
	case string:
		var categories []entity.KenaliDiriCategory
		if err := json.Unmarshal([]byte(v), &categories); err == nil {
			return categories, true
		}
	}

	return nil, false
}
