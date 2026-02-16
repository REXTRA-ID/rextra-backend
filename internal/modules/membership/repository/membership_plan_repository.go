package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipPlanRepository interface {
		Create(ctx context.Context, tx *gorm.DB, membershipPlan entity.MembershipPlans) (entity.MembershipPlans, error)
		GetByID(ctx context.Context, tx *gorm.DB, membershipPlanId uuid.UUID) (entity.MembershipPlans, error)
		GetByPlanName(ctx context.Context, tx *gorm.DB, planName string) (entity.MembershipPlans, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error)
		Update(ctx context.Context, tx *gorm.DB, membershipPlan entity.MembershipPlans) (entity.MembershipPlans, error)
		Delete(ctx context.Context, tx *gorm.DB, membershipPlanId uuid.UUID) (entity.MembershipPlans, error)
	}

	membershipPlanRepository struct {
		db *gorm.DB
	}
)

func NewMembershipPlanRepository(db *gorm.DB) MembershipPlanRepository {
	return &membershipPlanRepository{
		db: db,
	}
}

func (r *membershipPlanRepository) Create(ctx context.Context, tx *gorm.DB, membershipPlan entity.MembershipPlans) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&membershipPlan).Error; err != nil {
		return entity.MembershipPlans{}, err
	}

	return membershipPlan, nil
}

func (r *membershipPlanRepository) GetByID(ctx context.Context, tx *gorm.DB, membershipPlanId uuid.UUID) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipPlan entity.MembershipPlans
	if err := tx.WithContext(ctx).First(&membershipPlan, "id = ?", membershipPlanId).Error; err != nil {
		return entity.MembershipPlans{}, err
	}

	return membershipPlan, nil
}

func (r *membershipPlanRepository) GetByPlanName(ctx context.Context, tx *gorm.DB, planName string) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipPlan entity.MembershipPlans
	if err := tx.WithContext(ctx).First(&membershipPlan, "plan_name = ?", entity.PLANSTARTER).Error; err != nil {
		return entity.MembershipPlans{}, err
	}

	return membershipPlan, nil
}

func (r *membershipPlanRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipPlans []entity.MembershipPlans
	if err := tx.WithContext(ctx).Find(&membershipPlans).Error; err != nil {
		return []entity.MembershipPlans{}, err
	}

	return membershipPlans, nil
}

func (r *membershipPlanRepository) Update(ctx context.Context, tx *gorm.DB, membershipPlan entity.MembershipPlans) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&membershipPlan).Error; err != nil {
		return entity.MembershipPlans{}, err
	}

	return membershipPlan, nil
}

func (r *membershipPlanRepository) Delete(ctx context.Context, tx *gorm.DB, membershipPlanId uuid.UUID) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipPlan entity.MembershipPlans
	if err := tx.WithContext(ctx).First(&membershipPlan, "id = ?", membershipPlanId).Delete(&membershipPlan, "id = ?", membershipPlanId).Error; err != nil {
		return entity.MembershipPlans{}, err
	}

	return membershipPlan, nil
}
