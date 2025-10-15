package repository

import (
	"context"
	"errors"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	EducationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, userEducation entity.Education) (entity.Education, error)
		GetActiveEducationByUserId(ctx context.Context, tx *gorm.DB, userId string) (entity.Education, bool,error)
		GetAllEducationByUserId(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Education, error)
		GetByUserIdAndEducationById(ctx context.Context, tx *gorm.DB, userID string, educationID string) (entity.Education, bool,error)
		Update(ctx context.Context, tx *gorm.DB, userEducation entity.Education) (entity.Education, error)
		Delete(ctx context.Context, tx *gorm.DB, userID string, educationID string) error
	}

	educationRepository struct {
		db *gorm.DB
	}
)

func NewEducation(db *gorm.DB) EducationRepository {
	return &educationRepository{db: db}
}

func (r *educationRepository) Create(ctx context.Context, tx *gorm.DB, userEducation entity.Education) (entity.Education, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&userEducation).Error; err != nil {
		return userEducation, err
	}

	return userEducation, nil
}

func (r *educationRepository) Update(ctx context.Context, tx *gorm.DB, userEducation entity.Education) (entity.Education, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&userEducation).Error; err != nil {
		return userEducation, err
	}

	return userEducation, nil
}


func (r *educationRepository) Delete(ctx context.Context, tx *gorm.DB, userID string, educationID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&entity.Education{}, "user_id = ? AND id = ?", userID, educationID).Error; err != nil {
		return err
	}

	return nil
}

func (r *educationRepository) GetAllByUserId(ctx context.Context, tx *gorm.DB, userId uint) ([]entity.Education, error) {
	if tx == nil {
		tx = r.db
	}

	var userEducation []entity.Education
	if err := tx.WithContext(ctx).Where("user_id = ?", userId).Order("created_at DESC").Find(&userEducation).Error; err != nil {
		return userEducation, err
	}

	return userEducation, nil
}

func (r *educationRepository) GetByUserIdAndEducationById(ctx context.Context, tx *gorm.DB, userID string, educationID string) (entity.Education, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var userEducation entity.Education
	if err := tx.WithContext(ctx).Take(&userEducation, "user_id = ? AND id = ?", userID, educationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEducation, false, nil
		}
		return userEducation, false, err
	}

	return userEducation, true, nil
}

func (r *educationRepository) GetActiveEducationByUserId(ctx context.Context, tx *gorm.DB, userId string) (entity.Education, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var userEducation entity.Education
	if err := tx.WithContext(ctx).Take(&userEducation, "user_id = ? AND is_active = true", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEducation, false, nil
		}
		return userEducation, false, err
	}

	return userEducation, true, nil
}

func (r *educationRepository) GetAllEducationByUserId(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Education, error) {
	if tx == nil {
		tx = r.db
	}

	var userEducation []entity.Education
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&userEducation).Error; err != nil {
		return userEducation, err
	}

	return userEducation, nil
}
