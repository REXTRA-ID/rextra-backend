package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipPlanRepository interface {
		Create(ctx context.Context, tx *gorm.DB, model entity.MembershipPlans) (entity.MembershipPlans, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error)
		GetAllWithDurations(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error)
		GetByPlanName(ctx context.Context, tx *gorm.DB, planName entity.EnumPlanName) (entity.MembershipPlans, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.MembershipPlans) (entity.MembershipPlans, error)
		Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
		CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error)
	}

	membershipPlanRepository struct {
		db *gorm.DB
	}
)

func NewMembershipPlanRepository(db *gorm.DB) MembershipPlanRepository {
	return &membershipPlanRepository{db: db}
}

func (r *membershipPlanRepository) Create(ctx context.Context, tx *gorm.DB, model entity.MembershipPlans) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil {
		return entity.MembershipPlans{}, err
	}
	return model, nil
}

func (r *membershipPlanRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.MembershipPlans
	if err := tx.WithContext(ctx).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// GetAllWithDurations mengambil semua plan beserta PlanDurations-nya (preload).
// Dipakai untuk halaman katalog admin yang perlu tampilkan harga per durasi sekaligus.
func (r *membershipPlanRepository) GetAllWithDurations(ctx context.Context, tx *gorm.DB) ([]entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	var models []entity.MembershipPlans
	if err := tx.WithContext(ctx).
		Preload("PlanDurations", "is_active = ?", true).
		Order("created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *membershipPlanRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.MembershipPlans
	if err := tx.WithContext(ctx).
		Preload("PlanDurations", func(db *gorm.DB) *gorm.DB {
			return db.Order("duration_months ASC")
		}).
		First(&model, "id = ?", id).Error; err != nil {
		return entity.MembershipPlans{}, err
	}
	return model, nil
}

func (r *membershipPlanRepository) GetByPlanName(ctx context.Context, tx *gorm.DB, planName entity.EnumPlanName) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	var model entity.MembershipPlans
	if err := tx.WithContext(ctx).First(&model, "plan_name = ?", planName).Error; err != nil {
		return entity.MembershipPlans{}, err
	}
	return model, nil
}

func (r *membershipPlanRepository) Update(ctx context.Context, tx *gorm.DB, model entity.MembershipPlans) (entity.MembershipPlans, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil {
		return entity.MembershipPlans{}, err
	}
	return model, nil
}

func (r *membershipPlanRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.WithContext(ctx).Delete(&entity.MembershipPlans{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// CountActiveMembers menghitung jumlah Memberships aktif yang menggunakan plan ini.
// Guard sebelum Delete plan dan sebelum ubah PlanName/Category/DurationMode.
func (r *membershipPlanRepository) CountActiveMembers(ctx context.Context, tx *gorm.DB, id uuid.UUID) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	if err := tx.WithContext(ctx).Model(&entity.Memberships{}).
		Where("plan_id = ? AND is_active = true", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
