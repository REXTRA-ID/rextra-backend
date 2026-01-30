package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	RecommendationRepository interface {
		Save(ctx context.Context, tx *gorm.DB, recommendation entity.CareerRecommendation) error
		GetBySessionID(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.CareerRecommendation, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.CareerRecommendation, error)
	}

	recommendationRepository struct {
		db *gorm.DB
	}
)

func NewRecommendation(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Save(ctx context.Context, tx *gorm.DB, recommendation entity.CareerRecommendation) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&recommendation).Error
}

func (r *recommendationRepository) GetBySessionID(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	var recommendation entity.CareerRecommendation
	if err := tx.WithContext(ctx).Take(&recommendation, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.CareerRecommendation{}, err
	}

	return recommendation, nil
}

func (r *recommendationRepository) GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	var recommendation entity.CareerRecommendation
	if err := tx.WithContext(ctx).Take(&recommendation, "id = ?", id).Error; err != nil {
		return entity.CareerRecommendation{}, err
	}

	return recommendation, nil
}
