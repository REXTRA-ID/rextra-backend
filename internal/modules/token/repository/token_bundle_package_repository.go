package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	TokenBundlePackageRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenBundlePackage, error)
		Create(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error)
		Update(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) error
		ToggleActive(ctx context.Context, tx *gorm.DB, id string, isActive bool) error
	}
	tokenBundlePackageRepository struct {
		db *gorm.DB
	}
)

func NewTokenBundlePackageRepository(db *gorm.DB) TokenBundlePackageRepository {
	return &tokenBundlePackageRepository{db: db}
}

func (r *tokenBundlePackageRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error) {
	if tx == nil {
		tx = r.db
	}
	var tokenBundlePackages []entity.TokenBundlePackage
	if err := tx.Find(&tokenBundlePackages).Error; err != nil {
		return nil, err
	}
	return tokenBundlePackages, nil
}

func (r *tokenBundlePackageRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenBundlePackage, error) {
	if tx == nil {
		tx = r.db
	}
	var tokenBundlePackage entity.TokenBundlePackage
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&tokenBundlePackage).Error; err != nil {
		return entity.TokenBundlePackage{}, err
	}
	return tokenBundlePackage, nil
}

func (r *tokenBundlePackageRepository) Create(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Create(tokenBundlePackage).Error; err != nil {
		return entity.TokenBundlePackage{}, err
	}
	return tokenBundlePackage, nil
}

func (r *tokenBundlePackageRepository) Update(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Save(tokenBundlePackage).Error; err != nil {
		return err
	}
	return nil
}

func (r *tokenBundlePackageRepository) ToggleActive(ctx context.Context, tx *gorm.DB, id string, isActive bool) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Model(&entity.TokenBundlePackage{}).Where("id = ?", id).Update("active", isActive).Error; err != nil {
		return err
	}
	return nil
}
