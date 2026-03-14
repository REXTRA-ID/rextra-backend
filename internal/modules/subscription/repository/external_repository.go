package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByIDs(ctx context.Context, tx *gorm.DB, ids []uuid.UUID) (map[uuid.UUID]entity.User, error)
}
