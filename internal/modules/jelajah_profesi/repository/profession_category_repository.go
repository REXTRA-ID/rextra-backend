package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	ProfessionCategoryRepository interface {
		ListMainCategories(ctx context.Context) ([]entity.ProfessionMainCategory, error)
		ListSubCategories(ctx context.Context, mainCategoryID int64) ([]entity.ProfessionSubCategory, error)
	}

	professionCategoryRepository struct {
		db *gorm.DB
	}
)

func NewProfessionCategory(db *gorm.DB) ProfessionCategoryRepository {
	return &professionCategoryRepository{db: db}
}

func (r *professionCategoryRepository) ListMainCategories(ctx context.Context) ([]entity.ProfessionMainCategory, error) {
	var categories []entity.ProfessionMainCategory
	if err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *professionCategoryRepository) ListSubCategories(ctx context.Context, mainCategoryID int64) ([]entity.ProfessionSubCategory, error) {
	var categories []entity.ProfessionSubCategory
	if err := r.db.WithContext(ctx).
		Where("main_category_id = ?", mainCategoryID).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
