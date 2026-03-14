package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	CareerRecommendationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.CareerRecommendation) (entity.CareerRecommendation, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.CareerRecommendation, error)
		GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.CareerRecommendation, error)
		Update(ctx context.Context, tx *gorm.DB, session entity.CareerRecommendation) (entity.CareerRecommendation, error)
	}

	careerRecommendationRepository struct {
		db *gorm.DB
	}
)

func NewCareerRecommendation(db *gorm.DB) CareerRecommendationRepository {
	return &careerRecommendationRepository{db}
}

func (r *careerRecommendationRepository) Create(ctx context.Context, tx *gorm.DB, careerRecommendation entity.CareerRecommendation) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&careerRecommendation).Error; err != nil {
		return careerRecommendation, err
	}

	return careerRecommendation, nil
}

func (r *careerRecommendationRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	var careerRecommendation entity.CareerRecommendation
	if err := tx.WithContext(ctx).Take(&careerRecommendation, "id = ?", id).Error; err != nil {
		return entity.CareerRecommendation{}, err
	}

	return careerRecommendation, nil
}

func (r *careerRecommendationRepository) GetByUserID(ctx context.Context, tx *gorm.DB, userID string) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	var careerRecommendation entity.CareerRecommendation
	if err := tx.WithContext(ctx).Take(&careerRecommendation, "user_id = ?", userID).Error; err != nil {
		return entity.CareerRecommendation{}, err
	}
	return careerRecommendation, nil
}

func (r *careerRecommendationRepository) Update(ctx context.Context, tx *gorm.DB, careerRecommendation entity.CareerRecommendation) (entity.CareerRecommendation, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&careerRecommendation).Error; err != nil {
		return careerRecommendation, err
	}

	return careerRecommendation, nil
}
