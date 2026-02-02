package repository

import (
	"context"
	"errors"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	TokenBundlePackageRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.TokenBundlePackage, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.TokenBundlePackage, error)
		GetByName(ctx context.Context, tx *gorm.DB, name string) (entity.TokenBundlePackage, bool, error)
		Create(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error)
		Update(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error)
		Delete(ctx context.Context, tx *gorm.DB, id string) error
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

func (r *tokenBundlePackageRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (entity.TokenBundlePackage, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var tokenBundlePackage entity.TokenBundlePackage
	if err := tx.WithContext(ctx).Where("name = ?", name).First(&tokenBundlePackage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.TokenBundlePackage{}, false, nil
		}
		return entity.TokenBundlePackage{}, false, err
	}
	return tokenBundlePackage, true, nil
}

func (r *tokenBundlePackageRepository) Create(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Create(&tokenBundlePackage).Error; err != nil {
		return entity.TokenBundlePackage{}, err
	}
	return tokenBundlePackage, nil
}

func (r *tokenBundlePackageRepository) Update(ctx context.Context, tx *gorm.DB, tokenBundlePackage entity.TokenBundlePackage) (entity.TokenBundlePackage, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Save(&tokenBundlePackage).Error; err != nil {
		return entity.TokenBundlePackage{}, err
	}
	return tokenBundlePackage, nil
}

func (r *tokenBundlePackageRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = r.db
	}
	result := tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.TokenBundlePackage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return myerror.RecordNotFound("bundle")
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
