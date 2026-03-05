package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	ProfessionDetailRepository interface {
		GetActivities(ctx context.Context, professionID int64) ([]entity.ProfessionActivity, error)
		GetSkills(ctx context.Context, professionID int64) ([]entity.ProfessionSkill, error)
		GetTools(ctx context.Context, professionID int64) ([]entity.ProfessionTool, error)
		GetCareerPaths(ctx context.Context, professionID int64) ([]entity.ProfessionCareerPath, error)
		GetMarketInsights(ctx context.Context, professionID int64) ([]entity.ProfessionMarketInsight, error)
		GetStudyPrograms(ctx context.Context, professionID int64) ([]entity.ProfessionStudyProgram, error)
		GetAliases(ctx context.Context, professionID int64) ([]entity.ProfessionAlias, error)
	}

	professionDetailRepository struct {
		db *gorm.DB
	}
)

func NewProfessionDetail(db *gorm.DB) ProfessionDetailRepository {
	return &professionDetailRepository{db: db}
}

func (r *professionDetailRepository) GetActivities(ctx context.Context, professionID int64) ([]entity.ProfessionActivity, error) {
	var items []entity.ProfessionActivity
	if err := r.db.WithContext(ctx).
		Where("profession_id = ?", professionID).
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetSkills(ctx context.Context, professionID int64) ([]entity.ProfessionSkill, error) {
	var items []entity.ProfessionSkill
	if err := r.db.WithContext(ctx).
		Preload("Skill").
		Where("profession_id = ?", professionID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetTools(ctx context.Context, professionID int64) ([]entity.ProfessionTool, error) {
	var items []entity.ProfessionTool
	if err := r.db.WithContext(ctx).
		Preload("Tool").
		Where("profession_id = ?", professionID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetCareerPaths(ctx context.Context, professionID int64) ([]entity.ProfessionCareerPath, error) {
	var items []entity.ProfessionCareerPath
	if err := r.db.WithContext(ctx).
		Where("profession_id = ?", professionID).
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetMarketInsights(ctx context.Context, professionID int64) ([]entity.ProfessionMarketInsight, error) {
	var items []entity.ProfessionMarketInsight
	if err := r.db.WithContext(ctx).
		Where("profession_id = ?", professionID).
		Order("sort_order ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetStudyPrograms(ctx context.Context, professionID int64) ([]entity.ProfessionStudyProgram, error) {
	var items []entity.ProfessionStudyProgram
	if err := r.db.WithContext(ctx).
		Preload("StudyProgram").
		Where("profession_id = ?", professionID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *professionDetailRepository) GetAliases(ctx context.Context, professionID int64) ([]entity.ProfessionAlias, error) {
	var items []entity.ProfessionAlias
	if err := r.db.WithContext(ctx).
		Where("profession_id = ?", professionID).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
