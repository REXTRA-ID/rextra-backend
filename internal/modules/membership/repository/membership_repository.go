package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipRepository interface {
		Create(ctx context.Context, tx *gorm.DB, membership entity.Memberships) (entity.Memberships, error)
		GetByID(ctx context.Context, tx *gorm.DB, membershipId uuid.UUID) (entity.Memberships, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entity.Memberships, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error)
		GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.Memberships, int, error)
		GetExpiredMemberships(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error)
		GetActiveMemberships(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error)
		Update(ctx context.Context, tx *gorm.DB, membership entity.Memberships) (entity.Memberships, error)
	}
	membershipRepository struct {
		db *gorm.DB
	}
)

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{
		db: db,
	}
}

func (r *membershipRepository) Create(ctx context.Context, tx *gorm.DB, membership entity.Memberships) (entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&membership).Error; err != nil {
		return entity.Memberships{}, err
	}

	return membership, nil
}

func (r *membershipRepository) GetByID(ctx context.Context, tx *gorm.DB, membershipId uuid.UUID) (entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	var membership entity.Memberships
	if err := tx.WithContext(ctx).First(&membership, "id = ?", membershipId).Error; err != nil {
		return entity.Memberships{}, err
	}

	return membership, nil
}

func (r *membershipRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId uuid.UUID) (entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	var userMembership entity.Memberships
	if err := tx.WithContext(ctx).First(&userMembership, "user_id = ?", userId).Error; err != nil {
		return entity.Memberships{}, err
	}

	return userMembership, nil
}

func (r *membershipRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	var memberships []entity.Memberships
	if err := tx.WithContext(ctx).Find(&memberships).Error; err != nil {
		return []entity.Memberships{}, err
	}

	return memberships, nil
}

func (r *membershipRepository) GetAllPaginated(ctx context.Context, tx *gorm.DB, page, totalPage int) ([]entity.Memberships, int, error) {
	if tx == nil {
		tx = r.db
	}

	offset := (page + 1) * totalPage

	var memberships []entity.Memberships
	if err := tx.WithContext(ctx).Find(&memberships).Offset(offset).Limit(totalPage).Error; err != nil {
		return []entity.Memberships{}, 0, err
	}

	var total int64
	if err := tx.WithContext(ctx).Model(&entity.Memberships{}).Count(&total).Error; err != nil {
		return []entity.Memberships{}, 0, err
	}

	return memberships, int(total), nil
}

func (r *membershipRepository) GetExpiredMemberships(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	var expiredMemberships []entity.Memberships
	if err := tx.WithContext(ctx).Where("expired_at < (NOW() AT TIME ZONE 'UTC')").Find(&expiredMemberships).Error; err != nil {
		return []entity.Memberships{}, err
	}

	return expiredMemberships, nil
}

func (r *membershipRepository) GetActiveMemberships(ctx context.Context, tx *gorm.DB) ([]entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	var expiredMemberships []entity.Memberships
	if err := tx.WithContext(ctx).Preload("Plan").Where("is_active = ?", true).Find(&expiredMemberships).Error; err != nil {
		return []entity.Memberships{}, err
	}

	return expiredMemberships, nil
}

func (r *membershipRepository) Update(ctx context.Context, tx *gorm.DB, membership entity.Memberships) (entity.Memberships, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&membership).Error; err != nil {
		return entity.Memberships{}, err
	}

	return membership, nil
}
