package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipDurationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.MembershipDuration) (entity.MembershipDuration, error)
		GetAllByPlanID(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.MembershipDuration, error)
		GetAllByPlanIDWithMappingCounts(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.MembershipDuration, map[uuid.UUID]int64, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipDuration, error)
		GetByPlanAndMonths(ctx context.Context, tx *gorm.DB, planID uuid.UUID, months int) (entity.MembershipDuration, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.MembershipDuration) (entity.MembershipDuration, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
		CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	membershipDurationRepository struct {
		db *gorm.DB
	}
)

func NewMembershipDurationRepository(db *gorm.DB) MembershipDurationRepository {
	return &membershipDurationRepository{db: db}
}

func (r *membershipDurationRepository) Create(ctx context.Context, tx *gorm.DB, model entity.MembershipDuration) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.MembershipDuration{}, err
	}
	return model, nil
}

func (r *membershipDurationRepository) GetAllByPlanID(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.MembershipDuration
	if err := tx.WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("duration_months ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *membershipDurationRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.MembershipDuration
	if err := tx.WithContext(ctx).
		Preload("Plan").
		First(&model, "id = ?", id).Error; err != nil {
		return entity.MembershipDuration{}, err
	}
	return model, nil
}

// GetAllByPlanIDWithMappingCounts mengambil semua MembershipDuration milik satu plan
// beserta mapping count masing-masing dalam dua query flat — menghindari N+1
func (r *membershipDurationRepository) GetAllByPlanIDWithMappingCounts(ctx context.Context, tx *gorm.DB, planID uuid.UUID) ([]entity.MembershipDuration, map[uuid.UUID]int64, error) {
	if tx == nil {
		tx = r.db
	}

	var models []entity.MembershipDuration
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

// dipakai untuk cek duplikasi sebelum Create
func (r *membershipDurationRepository) GetByPlanAndMonths(ctx context.Context, tx *gorm.DB, planID uuid.UUID, months int) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	var model entity.MembershipDuration
	if err := tx.WithContext(ctx).
		First(&model, "plan_id = ? AND duration_months = ?", planID, months).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return model, nil
}

func (r *membershipDurationRepository) Update(ctx context.Context, tx *gorm.DB, model entity.MembershipDuration) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return model, nil
}

func (r *membershipDurationRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.MembershipDuration{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *membershipDurationRepository) CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
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

func (r *membershipDurationRepository) CountMappings(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
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
