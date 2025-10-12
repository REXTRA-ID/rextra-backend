package repository

import (
	"context"

	"gorm.io/gorm"

	"rextra-backend/internal/entity"
)

type (
	AssesmentRepository interface {
		CreateUserRiasec(ctx context.Context, tx *gorm.DB, userRiasec entity.UserRiasec) (entity.UserRiasec, error)
		CreateUserIkigai(ctx context.Context, tx *gorm.DB, userIkigai entity.UserIkigai) (entity.UserIkigai, error)
		GetUserRiasec(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserRiasec, error)
		GetUserIkigai(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserIkigai, error)
	}

	assesmentRepository struct {
		db *gorm.DB
	}
)

func NewAssesmentRepository(db *gorm.DB) AssesmentRepository {
	return &assesmentRepository{db: db}
}

func (r *assesmentRepository) CreateUserRiasec(ctx context.Context, tx *gorm.DB, userRiasec entity.UserRiasec) (entity.UserRiasec, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&userRiasec).Error; err != nil {
		return userRiasec, err
	}

	return userRiasec, nil
}

func (r *assesmentRepository) CreateUserIkigai(ctx context.Context, tx *gorm.DB, userIkigai entity.UserIkigai) (entity.UserIkigai, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&userIkigai).Error; err != nil {
		return userIkigai, err
	}

	return userIkigai, nil
}

func (r *assesmentRepository) GetUserRiasec(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserRiasec, error) {
	if tx == nil {
		tx = r.db
	}

	var userRiasec []entity.UserRiasec
	if err := tx.WithContext(ctx).Take(&userRiasec, "user_id = ?", userID).Error; err != nil {
		return []entity.UserRiasec{}, err
	}

	return userRiasec, nil
}

func (r *assesmentRepository) GetUserIkigai(ctx context.Context, tx *gorm.DB, userID string) ([]entity.UserIkigai, error) {
	if tx == nil {
		tx = r.db
	}

	var userIkigai []entity.UserIkigai
	if err := tx.WithContext(ctx).Take(&userIkigai, "user_id = ?", userID).Error; err != nil {
		return []entity.UserIkigai{}, err
	}

	return userIkigai, nil
}
