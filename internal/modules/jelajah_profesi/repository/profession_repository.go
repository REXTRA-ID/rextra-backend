package repository

import (
	"context"
	"strings"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type ProfessionFilters struct {
	Search         string
	MainCategoryID *int64
	SubCategoryID  *int64
}

type (
	ProfessionRepository interface {
		List(ctx context.Context, filters ProfessionFilters, limit, offset int) ([]entity.Profession, int64, error)
		GetBySlug(ctx context.Context, slug string) (entity.Profession, error)
	}

	professionRepository struct {
		db *gorm.DB
	}
)

func NewProfession(db *gorm.DB) ProfessionRepository {
	return &professionRepository{db: db}
}

func (r *professionRepository) List(ctx context.Context, filters ProfessionFilters, limit, offset int) ([]entity.Profession, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Profession{})

	query = applyProfessionFilters(query, filters)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var professions []entity.Profession
	if err := query.
		Preload("SubCategory").
		Order("professions.name ASC").
		Limit(limit).
		Offset(offset).
		Find(&professions).Error; err != nil {
		return nil, 0, err
	}

	return professions, total, nil
}

func (r *professionRepository) GetBySlug(ctx context.Context, slug string) (entity.Profession, error) {
	var profession entity.Profession
	if err := r.db.WithContext(ctx).
		Preload("MainCategory").
		Preload("SubCategory").
		Preload("RiasecCode").
		Where("slug = ?", slug).
		Take(&profession).Error; err != nil {
		return entity.Profession{}, err
	}
	return profession, nil
}

func applyProfessionFilters(query *gorm.DB, filters ProfessionFilters) *gorm.DB {
	if filters.MainCategoryID != nil {
		query = query.Where("professions.main_category_id = ?", *filters.MainCategoryID)
	}
	if filters.SubCategoryID != nil {
		query = query.Where("professions.sub_category_id = ?", *filters.SubCategoryID)
	}
	if filters.Search != "" {
		keyword := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where(
			"LOWER(professions.name) LIKE ? OR professions.id IN (SELECT profession_id FROM profession_aliases WHERE LOWER(alias_name) LIKE ?)",
			keyword, keyword,
		)
	}
	return query
}
