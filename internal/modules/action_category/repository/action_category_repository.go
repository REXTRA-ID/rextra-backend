package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ActionCategoryRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.ActionCategory, error)
		GetAllWithEntitlementCount(ctx context.Context, tx *gorm.DB) ([]entity.ActionCategory, map[uuid.UUID]int64, error)
		GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.ActionCategory, error)
		GetBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.ActionCategory, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	actionCategoryRepository struct {
		db *gorm.DB
	}
)

func NewActionCategoryRepository(db *gorm.DB) ActionCategoryRepository {
	return &actionCategoryRepository{db: db}
}

func (r *actionCategoryRepository) Create(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
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

func (r *actionCategoryRepository) GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.ActionCategory
	if err := tx.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return entity.ActionCategory{}, err
	}
	return model, nil
}

func (r *actionCategoryRepository) GetBySlug(ctx context.Context, tx *gorm.DB, slug string) (entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.ActionCategory
	if err := tx.WithContext(ctx).First(&model, "slug = ?", slug).Error; err != nil {
		return entity.ActionCategory{}, err
	}
	return model, nil
}

func (r *actionCategoryRepository) Update(ctx context.Context, tx *gorm.DB, model entity.ActionCategory) (entity.ActionCategory, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.ActionCategory{}, err
	}
	return model, nil
}

func (r *actionCategoryRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.ActionCategory{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CountEntitlements menghitung berapa Entitlement yang menggunakan ActionCategory ini.
// Dipakai service untuk:
// (1) Cek apakah slug boleh diubah (immutable jika count > 0)
// (2) Cek apakah category boleh dihapus (tidak boleh jika count > 0)
func (r *actionCategoryRepository) CountEntitlements(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Where("action_category_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetAllWithEntitlementCount mengambil semua ActionCategory sekaligus dengan jumlah
// entitlement per category dalam satu query — menghindari N+1 query di GetAll service.
// Hasilnya adalah map[uuid.UUID]int64 yang bisa di-lookup saat mapping response.
func (r *actionCategoryRepository) GetAllWithEntitlementCount(ctx context.Context, tx *gorm.DB) ([]entity.ActionCategory, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.ActionCategory
	if err := tx.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, nil, err
	}

	// Satu query aggregate untuk semua count sekaligus
	type countResult struct {
		ActionCategoryID uuid.UUID `gorm:"column:action_category_id"`
		Count            int64     `gorm:"column:count"`
	}
	var counts []countResult
	if err := tx.WithContext(ctx).Model(&entity.Entitlement{}).
		Select("action_category_id, COUNT(*) as count").
		Group("action_category_id").
		Scan(&counts).Error; err != nil {
		return nil, nil, err
	}

	countMap := make(map[uuid.UUID]int64)
	for _, c := range counts {
		countMap[c.ActionCategoryID] = c.Count
	}

	return models, countMap, nil
}
