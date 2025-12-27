package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"rextra-backend/internal/entity"
	"rextra-backend/internal/pkg/cache"

	"gorm.io/gorm"
)

const (
	riasecCacheKeyAll      = "riasec_codes:all"
	riasecCacheKeyByID     = "riasec_code:%d"
	riasecCacheKeyByCode   = "riasec_code_by_name:%s"
	riasecCacheKeyTypeList = "riasec_codes:type:%s"
)

type (
	RiasecCodeRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.RiasecCode, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.RiasecCode, error)
		GetByCode(ctx context.Context, tx *gorm.DB, code string) (entity.RiasecCode, error)
		ListByType(ctx context.Context, tx *gorm.DB, codeType string) ([]entity.RiasecCode, error)
		Search(ctx context.Context, tx *gorm.DB, keyword string) ([]entity.RiasecCode, error)
		Update(ctx context.Context, tx *gorm.DB, code entity.RiasecCode) error
	}

	riasecCodeRepository struct {
		db    *gorm.DB
		cache cache.CacheService
		ttl   time.Duration
	}
)

func NewRiasecCode(db *gorm.DB, cacheSvc cache.CacheService) RiasecCodeRepository {
	return &riasecCodeRepository{
		db:    db,
		cache: cacheSvc,
		ttl:   24 * time.Hour,
	}
}

func (r *riasecCodeRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.RiasecCode, error) {
	if tx == nil {
		tx = r.db
	}

	if r.cache != nil {
		if cached, err := r.cache.Get(ctx, riasecCacheKeyAll); err == nil {
			if codes, ok := decodeRiasecCache(cached); ok {
				return codes, nil
			}
		}
	}

	var codes []entity.RiasecCode
	if err := tx.WithContext(ctx).Find(&codes).Error; err != nil {
		return nil, err
	}

	r.storeCache(ctx, riasecCacheKeyAll, codes)
	return codes, nil
}

func (r *riasecCodeRepository) GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.RiasecCode, error) {
	if tx == nil {
		tx = r.db
	}

	cacheKey := fmt.Sprintf(riasecCacheKeyByID, id)
	if r.cache != nil {
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if codes, ok := decodeRiasecCache(cached); ok && len(codes) > 0 {
				return codes[0], nil
			}
		}
	}

	var code entity.RiasecCode
	if err := tx.WithContext(ctx).Take(&code, "id = ?", id).Error; err != nil {
		return entity.RiasecCode{}, err
	}

	r.storeCache(ctx, cacheKey, []entity.RiasecCode{code})
	return code, nil
}

func (r *riasecCodeRepository) GetByCode(ctx context.Context, tx *gorm.DB, codeName string) (entity.RiasecCode, error) {
	if tx == nil {
		tx = r.db
	}

	cacheKey := fmt.Sprintf(riasecCacheKeyByCode, strings.ToUpper(codeName))
	if r.cache != nil {
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if codes, ok := decodeRiasecCache(cached); ok && len(codes) > 0 {
				return codes[0], nil
			}
		}
	}

	var code entity.RiasecCode
	if err := tx.WithContext(ctx).Take(&code, "riasec_code = ?", strings.ToUpper(codeName)).Error; err != nil {
		return entity.RiasecCode{}, err
	}

	r.storeCache(ctx, cacheKey, []entity.RiasecCode{code})
	return code, nil
}

func (r *riasecCodeRepository) ListByType(ctx context.Context, tx *gorm.DB, codeType string) ([]entity.RiasecCode, error) {
	if tx == nil {
		tx = r.db
	}

	cacheKey := fmt.Sprintf(riasecCacheKeyTypeList, codeType)
	if r.cache != nil {
		if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
			if codes, ok := decodeRiasecCache(cached); ok {
				return codes, nil
			}
		}
	}

	query := tx.WithContext(ctx).Model(&entity.RiasecCode{})
	switch codeType {
	case "single":
		query = query.Where("char_length(riasec_code) = 1")
	case "dual":
		query = query.Where("char_length(riasec_code) = 2")
	case "triple":
		query = query.Where("char_length(riasec_code) = 3")
	}

	var codes []entity.RiasecCode
	if err := query.Find(&codes).Error; err != nil {
		return nil, err
	}

	r.storeCache(ctx, cacheKey, codes)
	return codes, nil
}

func (r *riasecCodeRepository) Search(ctx context.Context, tx *gorm.DB, keyword string) ([]entity.RiasecCode, error) {
	if tx == nil {
		tx = r.db
	}

	var codes []entity.RiasecCode
	key := strings.ToLower(keyword)
	if err := tx.WithContext(ctx).
		Where("LOWER(riasec_title) LIKE ? OR LOWER(riasec_code) LIKE ?", "%"+key+"%", "%"+key+"%").
		Find(&codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

func (r *riasecCodeRepository) Update(ctx context.Context, tx *gorm.DB, code entity.RiasecCode) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&code).Error; err != nil {
		return err
	}

	r.invalidateCache(ctx, code.ID, code.RiasecCode)
	return nil
}

func (r *riasecCodeRepository) storeCache(ctx context.Context, key string, data []entity.RiasecCode) {
	if r.cache == nil {
		return
	}

	if bytes, err := json.Marshal(data); err == nil {
		_ = r.cache.Set(ctx, key, bytes, r.ttl)
	}
}

func (r *riasecCodeRepository) invalidateCache(ctx context.Context, id int64, code string) {
	if r.cache == nil {
		return
	}

	_ = r.cache.Delete(ctx, fmt.Sprintf(riasecCacheKeyByID, id))
	_ = r.cache.Delete(ctx, fmt.Sprintf(riasecCacheKeyByCode, strings.ToUpper(code)))
	_ = r.cache.Delete(ctx, riasecCacheKeyAll)
	_ = r.cache.Clear(ctx, "riasec_codes:type*")
}

func decodeRiasecCache(data interface{}) ([]entity.RiasecCode, bool) {
	switch v := data.(type) {
	case []entity.RiasecCode:
		return v, true
	case []byte:
		var codes []entity.RiasecCode
		if err := json.Unmarshal(v, &codes); err == nil {
			return codes, true
		}
	case string:
		var codes []entity.RiasecCode
		if err := json.Unmarshal([]byte(v), &codes); err == nil {
			return codes, true
		}
	}

	return nil, false
}
