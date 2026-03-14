package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EntitlementRepository interface {
	GetById(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.Entitlement, error)
}
