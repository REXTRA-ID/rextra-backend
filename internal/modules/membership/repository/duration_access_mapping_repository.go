package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DurationAccessMappingRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.DurationAccessMapping) (entity.DurationAccessMapping, error)
		GetAllByPlanDurationID(ctx context.Context, tx *gorm.DB, planDurationID uuid.UUID) ([]entity.DurationAccessMapping, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.DurationAccessMapping, error)
		GetByPlanDurationAndEntitlement(ctx context.Context, tx *gorm.DB, planDurationID, entitlementID uuid.UUID) (entity.DurationAccessMapping, error)
		GetActiveByPlanDurationID(ctx context.Context, tx *gorm.DB, planDurationID uuid.UUID) ([]entity.DurationAccessMapping, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.DurationAccessMapping) (entity.DurationAccessMapping, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	}

	durationAccessMappingRepository struct {
		db *gorm.DB
	}
)

func NewDurationAccessMappingRepository(db *gorm.DB) DurationAccessMappingRepository {
	return &durationAccessMappingRepository{db: db}
}

func (r *durationAccessMappingRepository) Create(ctx context.Context, tx *gorm.DB, model entity.DurationAccessMapping) (entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.DurationAccessMapping{}, err
	}
	return model, nil
}

func (r *durationAccessMappingRepository) GetAllByPlanDurationID(ctx context.Context, tx *gorm.DB, planDurationID uuid.UUID) ([]entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.DurationAccessMapping
	if err := tx.WithContext(ctx).
		Preload("Entitlement").
		Where("plan_duration_id = ?", planDurationID).
		Order("created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *durationAccessMappingRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.DurationAccessMapping
	if err := tx.WithContext(ctx).
		Preload("Entitlement").
		Preload("PlanDuration").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.DurationAccessMapping{}, err
	}
	return model, nil
}

func (r *durationAccessMappingRepository) GetByPlanDurationAndEntitlement(ctx context.Context, tx *gorm.DB, planDurationID, entitlementID uuid.UUID) (entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.DurationAccessMapping
	if err := tx.WithContext(ctx).
		First(&model, "plan_duration_id = ? AND entitlement_id = ?", planDurationID, entitlementID).Error; err != nil {
		return entity.DurationAccessMapping{}, err
	}
	return model, nil
}

// GetActiveByPlanDurationID dipakai oleh access check runtime —
// hanya load mapping yang statusnya aktif untuk performa maksimal.
func (r *durationAccessMappingRepository) GetActiveByPlanDurationID(ctx context.Context, tx *gorm.DB, planDurationID uuid.UUID) ([]entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.DurationAccessMapping
	if err := tx.WithContext(ctx).
		Where("plan_duration_id = ? AND status = ?", planDurationID, entity.MappingStatusActive).
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *durationAccessMappingRepository) Update(ctx context.Context, tx *gorm.DB, model entity.DurationAccessMapping) (entity.DurationAccessMapping, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.DurationAccessMapping{}, err
	}
	return model, nil
}

func (r *durationAccessMappingRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.DurationAccessMapping{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
