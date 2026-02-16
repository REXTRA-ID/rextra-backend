package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	RiasecRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.Riasec) (entity.Riasec, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Riasec, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.Riasec, error)
		Update(ctx context.Context, tx *gorm.DB, session entity.Riasec) (entity.Riasec, error)
	}

	riasecRepository struct {
		db *gorm.DB
	}
)

func NewRiasec(db *gorm.DB) RiasecRepository {
	return &riasecRepository{db}
}

func (r *riasecRepository) Create(ctx context.Context, tx *gorm.DB, riasec entity.Riasec) (entity.Riasec, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&riasec).Error; err != nil {
		return riasec, err
	}

	return riasec, nil
}

func (r *riasecRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.Riasec, error) {
	if tx == nil {
		tx = r.db
	}

	var riasec entity.Riasec
	if err := tx.WithContext(ctx).Take(&riasec, "id = ?", id).Error; err != nil {
		return entity.Riasec{}, err
	}

	return riasec, nil
}

func (r *riasecRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userId string) (entity.Riasec, error) {
	if tx == nil {
		tx = r.db
	}

	var riasec entity.Riasec
	if err := tx.WithContext(ctx).Take(&riasec, "user_id = ?", userId).Error; err != nil {
		return entity.Riasec{}, err
	}

	return riasec, nil
}

func (r *riasecRepository) Update(ctx context.Context, tx *gorm.DB, riasec entity.Riasec) (entity.Riasec, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&riasec).Error; err != nil {
		return entity.Riasec{}, err
	}

	return riasec, nil
}
