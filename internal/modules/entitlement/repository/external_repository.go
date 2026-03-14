package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeatureRepository interface {
	GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Feature, error)
}

type SubFeatureRepository interface {
	GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.SubFeature, error)
}

type ActionCategoryRepository interface {
	GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.ActionCategory, error)
}
