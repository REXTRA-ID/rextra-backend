package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	UserMembershipRepository interface {
		GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, planName, userSearch string, isActive *bool) ([]entity.Memberships, int64, error)
		GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Memberships, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error)
		Create(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error)
		Update(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error)
	}

	userMembershipRepository struct {
		db *gorm.DB
	}
)

func NewUserMembershipRepository(db *gorm.DB) UserMembershipRepository {
	return &userMembershipRepository{db: db}
}

func (r *userMembershipRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, offset, limit int, planName, userSearch string, isActive *bool) ([]entity.Memberships, int64, error) {
	if tx == nil { tx = r.db }
	query := tx.WithContext(ctx).Model(&entity.Memberships{}).Joins("JOIN users ON users.id = memberships.user_id").Preload("Plan").Preload("Duration")
	if planName != "" { query = query.Where("memberships.plan_name = ?", planName) }
	if isActive != nil { query = query.Where("memberships.is_active = ?", *isActive) }
	if userSearch != "" {
		pattern := "%" + userSearch + "%"
		query = query.Where("users.fullname ILIKE ? OR users.email ILIKE ?", pattern, pattern)
	}
	var total int64
	if err := query.Select("COUNT(DISTINCT memberships.id)").Count(&total).Error; err != nil { return nil, 0, err }
	var results []entity.Memberships
	if err := query.Select("memberships.*").Order("memberships.updated_at DESC").Offset(offset).Limit(limit).Find(&results).Error; err != nil { return nil, 0, err }
	return results, total, nil
}

func (r *userMembershipRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Memberships, error) {
	if tx == nil { tx = r.db }
	var result entity.Memberships
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").First(&result, "id = ?", id).Error; err != nil { return entity.Memberships{}, err }
	return result, nil
}

func (r *userMembershipRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (entity.Memberships, error) {
	if tx == nil { tx = r.db }
	var result entity.Memberships
	if err := tx.WithContext(ctx).Preload("Plan").Preload("Duration").First(&result, "user_id = ?", userID).Error; err != nil { return entity.Memberships{}, err }
	return result, nil
}

func (r *userMembershipRepository) Create(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Create(&model).Error; err != nil { return entity.Memberships{}, err }
	return model, nil
}

func (r *userMembershipRepository) Update(ctx context.Context, tx *gorm.DB, model entity.Memberships) (entity.Memberships, error) {
	if tx == nil { tx = r.db }
	if err := tx.WithContext(ctx).Save(&model).Error; err != nil { return entity.Memberships{}, err }
	return model, nil
}
