package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	SessionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.SessionToken) (entity.SessionToken, error)
		GetByToken(ctx context.Context, tx *gorm.DB, token string) (entity.SessionToken, error)
		Update(ctx context.Context, tx *gorm.DB, session entity.SessionToken) (entity.SessionToken, error)
		Delete(ctx context.Context, tx *gorm.DB, session entity.SessionToken) error
	}

	sessionRepository struct {
		db *gorm.DB
	}
)

func NewSession(db *gorm.DB) SessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) Create(ctx context.Context, tx *gorm.DB, session entity.SessionToken) (entity.SessionToken, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, tx *gorm.DB, token string) (entity.SessionToken, error) {
	if tx == nil {
		tx = r.db
	}

	var session entity.SessionToken
	if err := tx.WithContext(ctx).Take(&session, "token = ?", token).Error; err != nil {
		return entity.SessionToken{}, err
	}

	return session, nil
}

func (r *sessionRepository) Update(ctx context.Context, tx *gorm.DB, session entity.SessionToken) (entity.SessionToken, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&session).Error; err != nil {
		return entity.SessionToken{}, err
	}

	return session, nil
}

func (r *sessionRepository) Delete(ctx context.Context, tx *gorm.DB, session entity.SessionToken) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&session).Error; err != nil {
		return err
	}

	return nil
}
