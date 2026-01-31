package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	RiasecRepository interface {
		SaveQuestionSet(ctx context.Context, tx *gorm.DB, questionSet entity.RiasecQuestionSet) error
		GetQuestionSet(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecQuestionSet, error)
		SaveResponses(ctx context.Context, tx *gorm.DB, responses entity.RiasecResponse) error
		GetResponses(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecResponse, error)
		SaveResult(ctx context.Context, tx *gorm.DB, result entity.RiasecResult) error
		GetResultBySessionID(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecResult, error)
		GetResultByID(ctx context.Context, tx *gorm.DB, id int64) (entity.RiasecResult, error)
	}

	riasecRepository struct {
		db *gorm.DB
	}
)

func NewRiasec(db *gorm.DB) RiasecRepository {
	return &riasecRepository{db: db}
}

func (r *riasecRepository) SaveQuestionSet(ctx context.Context, tx *gorm.DB, questionSet entity.RiasecQuestionSet) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&questionSet).Error
}

func (r *riasecRepository) GetQuestionSet(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecQuestionSet, error) {
	if tx == nil {
		tx = r.db
	}

	var questionSet entity.RiasecQuestionSet
	if err := tx.WithContext(ctx).Take(&questionSet, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.RiasecQuestionSet{}, err
	}

	return questionSet, nil
}

func (r *riasecRepository) SaveResponses(ctx context.Context, tx *gorm.DB, responses entity.RiasecResponse) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&responses).Error
}

func (r *riasecRepository) GetResponses(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var responses entity.RiasecResponse
	if err := tx.WithContext(ctx).Take(&responses, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.RiasecResponse{}, err
	}

	return responses, nil
}

func (r *riasecRepository) SaveResult(ctx context.Context, tx *gorm.DB, result entity.RiasecResult) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&result).Error
}

func (r *riasecRepository) GetResultBySessionID(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.RiasecResult, error) {
	if tx == nil {
		tx = r.db
	}

	var result entity.RiasecResult
	if err := tx.WithContext(ctx).
		Preload("RiasecCode").
		Take(&result, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.RiasecResult{}, err
	}

	return result, nil
}

func (r *riasecRepository) GetResultByID(ctx context.Context, tx *gorm.DB, id int64) (entity.RiasecResult, error) {
	if tx == nil {
		tx = r.db
	}

	var result entity.RiasecResult
	if err := tx.WithContext(ctx).
		Preload("RiasecCode").
		Take(&result, "id = ?", id).Error; err != nil {
		return entity.RiasecResult{}, err
	}

	return result, nil
}
