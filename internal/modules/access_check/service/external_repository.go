package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository interface {
	GetWalletByUserID(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (TokenWalletSnapshot, error)
	DeductToken(ctx context.Context, tx *gorm.DB, userID uuid.UUID, amount int, description string) error
}

type TokenWalletSnapshot struct {
	UserID  uuid.UUID
	Balance int64
}
