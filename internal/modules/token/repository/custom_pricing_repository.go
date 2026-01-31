package repository

import (
	"context"
	"rextra-backend/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomPricingRepository interface {
	// Config operations
	GetCurrentConfig(ctx context.Context) (*entity.CustomPricingConfig, error)
	GetConfigByID(ctx context.Context, id uuid.UUID) (*entity.CustomPricingConfig, error)
	GetConfigHistory(ctx context.Context, limit int) ([]entity.CustomPricingConfig, error)

	// Version management
	CreateNewVersion(ctx context.Context, config *entity.CustomPricingConfig, tiers []entity.CustomPricingTier) error
}

type customPricingRepository struct {
	db *gorm.DB
}

func NewCustomPricingRepository(db *gorm.DB) CustomPricingRepository {
	return &customPricingRepository{db: db}
}

// GetCurrentConfig retrieves current active config
func (r *customPricingRepository) GetCurrentConfig(ctx context.Context) (*entity.CustomPricingConfig, error) {
	var config entity.CustomPricingConfig

	err := r.db.WithContext(ctx).
		Preload("Tiers", func(db *gorm.DB) *gorm.DB {
			return db.Order("from_token ASC")
		}).
		Where("is_current = ?", true).
		First(&config).Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *customPricingRepository) GetConfigByID(ctx context.Context, id uuid.UUID) (*entity.CustomPricingConfig, error) {
	var config entity.CustomPricingConfig

	err := r.db.WithContext(ctx).
		Preload("Tiers", func(db *gorm.DB) *gorm.DB {
			return db.Order("from_token ASC")
		}).
		Preload("UpdatedByUser").
		Preload("GuardrailAcker").
		Where("id = ?", id).
		First(&config).Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *customPricingRepository) GetConfigHistory(ctx context.Context, limit int) ([]entity.CustomPricingConfig, error) {
	var configs []entity.CustomPricingConfig

	err := r.db.WithContext(ctx).
		Order("effective_from DESC").
		Limit(limit).
		Find(&configs).Error

	return configs, err
}

func (r *customPricingRepository) CreateNewVersion(ctx context.Context, config *entity.CustomPricingConfig, tiers []entity.CustomPricingTier) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Close current config
		if err := tx.Model(&entity.CustomPricingConfig{}).
			Where("is_current = ?", true).
			Updates(map[string]any{
				"is_current":   false,
				"effective_to": time.Now(),
			}).Error; err != nil {
			return err
		}

		// 2. Create new config
		config.IsCurrent = true
		config.EffectiveFrom = time.Now()
		if err := tx.Create(config).Error; err != nil {
			return err
		}

		// 3. Create tiers
		for i := range tiers {
			tiers[i].ConfigID = config.ID
			if err := tx.Create(&tiers[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
