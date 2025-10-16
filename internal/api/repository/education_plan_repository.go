package repository

import (
	"context"
	"errors"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	EducationPlanRepository interface {
		Create(ctx context.Context, tx *gorm.DB, educationPlan entity.EducationPlan) (entity.EducationPlan, error)
		GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string) ([]entity.EducationPlan, error)
		GetByUserIdAndEducationPlanById(ctx context.Context, tx *gorm.DB, userID string, id string) (entity.EducationPlan, bool, error)
		Update(ctx context.Context, tx *gorm.DB, educationPlan entity.EducationPlan) (entity.EducationPlan, error)
		Delete(ctx context.Context, tx *gorm.DB, userID string, id string) error
	}

	educationPlanRepository struct {
		db *gorm.DB
	}
)

func NewEducationPlan(db *gorm.DB) EducationPlanRepository {
	return &educationPlanRepository{db: db}
}

func (r *educationPlanRepository) Create(ctx context.Context, tx *gorm.DB, educationPlan entity.EducationPlan) (entity.EducationPlan, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&educationPlan).Error; err != nil {
		return educationPlan, err
	}

	return educationPlan, nil
}

func (r *educationPlanRepository) GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string) ([]entity.EducationPlan, error) {
	if tx == nil {
		tx = r.db
	}

	var plans []entity.EducationPlan
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&plans).Error; err != nil {
		return plans, err
	}

	return plans, nil
}

func (r *educationPlanRepository) GetByUserIdAndEducationPlanById(ctx context.Context, tx *gorm.DB, userID string, id string) (entity.EducationPlan, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var plan entity.EducationPlan
	if err := tx.WithContext(ctx).Take(&plan, "user_id = ? AND id = ?", userID, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return plan, false, nil
		}
		return plan, false, err
	}

	return plan, true, nil
}

func (r *educationPlanRepository) Update(ctx context.Context, tx *gorm.DB, educationPlan entity.EducationPlan) (entity.EducationPlan, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&educationPlan).Error; err != nil {
		return educationPlan, err
	}

	return educationPlan, nil
}

func (r *educationPlanRepository) Delete(ctx context.Context, tx *gorm.DB, userID string, id string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.EducationPlan{}, "user_id = ? AND id = ?", userID, id).Error; err != nil {
		return err
	}

	return nil
}