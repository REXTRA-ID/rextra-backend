package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PlanDurationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.PlanDuration) (entity.PlanDuration, error)
		GetAllByPlanID(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.PlanDuration, error)
		GetAllByPlanIDWithMappingCounts(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.PlanDuration, map[uuid.UUID]int64, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PlanDuration, error)
		GetByPlanAndMonths(ctx context.Context, tx *gorm.DB, planID uuid.UUID, months int) (entity.PlanDuration, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.PlanDuration) (entity.PlanDuration, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
		CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	planDurationRepository struct {
		db *gorm.DB
	}
)

func NewPlanDurationRepository(db *gorm.DB) PlanDurationRepository {
	return &planDurationRepository{db: db}
}

func (r *planDurationRepository) Create(ctx context.Context, tx *gorm.DB, model entity.PlanDuration) (entity.PlanDuration, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.PlanDuration{}, err
	}
	return model, nil
}

func (r *planDurationRepository) GetAllByPlanID(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.PlanDuration, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.PlanDuration
	if err := tx.WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("duration_months ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *planDurationRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.PlanDuration, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.PlanDuration
	if err := tx.WithContext(ctx).
		Preload("Plan").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.PlanDuration{}, err
	}
	return model, nil
}

func (r *planDurationRepository) GetAllByPlanIDWithMappingCounts(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.PlanDuration, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.PlanDuration
	if err := tx.WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("duration_months ASC").
		Find(&models).Error; err != nil {
		return nil, nil, err
	}

	ids := make([]uuid.UUID, len(models))
	for i, m := range models {
		ids[i] = m.ID
	}

	type countResult struct {
		PlanDurationID uuid.UUID `gorm:"column:plan_duration_id"`
		Count          int64     `gorm:"column:count"`
	}
	var counts []countResult
	if len(ids) > 0 {
		if err := tx.WithContext(ctx).Model(&entity.DurationAccessMapping{}).
			Select("plan_duration_id, COUNT(*) as count").
			Where("plan_duration_id IN ?", ids).
			Group("plan_duration_id").
			Scan(&counts).Error; err != nil {
			return nil, nil, err
		}
	}

	countMap := make(map[uuid.UUID]int64)
	for _, c := range counts {
		countMap[c.PlanDurationID] = c.Count
	}

	return models, countMap, nil
}

func (r *planDurationRepository) GetByPlanAndMonths(ctx context.Context, tx *gorm.DB, planID uuid.UUID, months int) (entity.PlanDuration, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.PlanDuration
	if err := tx.WithContext(ctx).
		First(&model, "plan_id = ? AND duration_months = ?", planID, months).Error; err != nil {
		return entity.PlanDuration{}, err
	}
	return model, nil
}

func (r *planDurationRepository) Update(ctx context.Context, tx *gorm.DB, model entity.PlanDuration) (entity.PlanDuration, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.PlanDuration{}, err
	}
	return model, nil
}

func (r *planDurationRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.PlanDuration{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *planDurationRepository) CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.Memberships{}).
		Where("duration_id = ? AND is_active = true", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *planDurationRepository) CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.DurationAccessMapping{}).
		Where("plan_duration_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
