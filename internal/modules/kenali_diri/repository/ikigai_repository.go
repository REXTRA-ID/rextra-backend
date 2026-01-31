package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	IkigaiRepository interface {
		SaveCandidates(ctx context.Context, tx *gorm.DB, candidates entity.IkigaiCandidateProfession) error
		GetCandidates(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiCandidateProfession, error)
		SaveResponses(ctx context.Context, tx *gorm.DB, responses entity.IkigaiResponse) error
		GetResponses(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiResponse, error)
		SaveDimensionScores(ctx context.Context, tx *gorm.DB, scores entity.IkigaiDimensionScore) error
		SaveTotalScores(ctx context.Context, tx *gorm.DB, scores entity.IkigaiTotalScore) error
		GetTotalScores(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiTotalScore, error)
	}

	ikigaiRepository struct {
		db *gorm.DB
	}
)

func NewIkigai(db *gorm.DB) IkigaiRepository {
	return &ikigaiRepository{db: db}
}

func (r *ikigaiRepository) SaveCandidates(ctx context.Context, tx *gorm.DB, candidates entity.IkigaiCandidateProfession) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&candidates).Error
}

func (r *ikigaiRepository) GetCandidates(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiCandidateProfession, error) {
	if tx == nil {
		tx = r.db
	}

	var candidates entity.IkigaiCandidateProfession
	if err := tx.WithContext(ctx).Take(&candidates, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.IkigaiCandidateProfession{}, err
	}

	return candidates, nil
}

func (r *ikigaiRepository) SaveResponses(ctx context.Context, tx *gorm.DB, responses entity.IkigaiResponse) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&responses).Error
}

func (r *ikigaiRepository) GetResponses(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var responses entity.IkigaiResponse
	if err := tx.WithContext(ctx).Take(&responses, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.IkigaiResponse{}, err
	}

	return responses, nil
}

func (r *ikigaiRepository) SaveDimensionScores(ctx context.Context, tx *gorm.DB, scores entity.IkigaiDimensionScore) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&scores).Error
}

func (r *ikigaiRepository) SaveTotalScores(ctx context.Context, tx *gorm.DB, scores entity.IkigaiTotalScore) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&scores).Error
}

func (r *ikigaiRepository) GetTotalScores(ctx context.Context, tx *gorm.DB, sessionID int64) (entity.IkigaiTotalScore, error) {
	if tx == nil {
		tx = r.db
	}

	var total entity.IkigaiTotalScore
	if err := tx.WithContext(ctx).Take(&total, "test_session_id = ?", sessionID).Error; err != nil {
		return entity.IkigaiTotalScore{}, err
	}

	return total, nil
}
