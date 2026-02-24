package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipDurationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, membershipDuration entity.MembershipDuration) (entity.MembershipDuration, error)
		GetByDurationMonth(ctx context.Context, tx *gorm.DB, duration int) (entity.MembershipDuration, error)
		GetByID(ctx context.Context, tx *gorm.DB, membershipDurationId uuid.UUID) (entity.MembershipDuration, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipDuration, error)
		Update(ctx context.Context, tx *gorm.DB, membershipDuration entity.MembershipDuration) (entity.MembershipDuration, error)
		Delete(ctx context.Context, tx *gorm.DB, membershipDurationId uuid.UUID) (entity.MembershipDuration, error)
	}

	membershipDurationRepository struct {
		db *gorm.DB
	}
)

func NewMembershipDurationRepository(db *gorm.DB) MembershipDurationRepository {
	return &membershipDurationRepository{
		db: db,
	}
}

func (r *membershipDurationRepository) Create(ctx context.Context, tx *gorm.DB, membershipDuration entity.MembershipDuration) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&membershipDuration).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return membershipDuration, nil
}

func (r *membershipDurationRepository) GetByDurationMonth(ctx context.Context, tx *gorm.DB, duration int) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipDuration entity.MembershipDuration
	if err := tx.WithContext(ctx).First(&membershipDuration, "duration_month = ?", duration).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return membershipDuration, nil
}

func (r *membershipDurationRepository) GetByID(ctx context.Context, tx *gorm.DB, membershipDurationId uuid.UUID) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	var data entity.MembershipDuration

	if err := tx.WithContext(ctx).
		First(&data, "id = ?", membershipDurationId).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return data, nil
}

func (r *membershipDurationRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	var durations []entity.MembershipDuration

	if err := tx.WithContext(ctx).Find(&durations).Error; err != nil {
		return []entity.MembershipDuration{}, err
	}

	return durations, nil
}

func (r *membershipDurationRepository) Update(ctx context.Context, tx *gorm.DB, membershipDuration entity.MembershipDuration) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Save(membershipDuration).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return membershipDuration, nil
}

func (r *membershipDurationRepository) Delete(ctx context.Context, tx *gorm.DB, membershipDurationId uuid.UUID) (entity.MembershipDuration, error) {
	if tx == nil {
		tx = r.db
	}

	var membershipDuration entity.MembershipDuration

	if err := tx.WithContext(ctx).
		First(&membershipDuration, "id = ?", membershipDurationId).
		Delete(&entity.MembershipDuration{}).Error; err != nil {
		return entity.MembershipDuration{}, err
	}

	return membershipDuration, nil
}
