package repository

import (
	"context"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	TestSessionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.CareerProfileTestSession) (entity.CareerProfileTestSession, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.CareerProfileTestSession, error)
		GetByToken(ctx context.Context, tx *gorm.DB, token string) (entity.CareerProfileTestSession, error)
		UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, timestamps map[string]interface{}) error
		GetUserSessions(ctx context.Context, tx *gorm.DB, userID string) ([]entity.CareerProfileTestSession, error)
	}

	testSessionRepository struct {
		db *gorm.DB
	}
)

func NewTestSession(db *gorm.DB) TestSessionRepository {
	return &testSessionRepository{db: db}
}

func (r *testSessionRepository) Create(ctx context.Context, tx *gorm.DB, session entity.CareerProfileTestSession) (entity.CareerProfileTestSession, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}

func (r *testSessionRepository) GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.CareerProfileTestSession, error) {
	if tx == nil {
		tx = r.db
	}

	var session entity.CareerProfileTestSession
	if err := tx.WithContext(ctx).Take(&session, "id = ?", id).Error; err != nil {
		return entity.CareerProfileTestSession{}, err
	}

	return session, nil
}

func (r *testSessionRepository) GetByToken(ctx context.Context, tx *gorm.DB, token string) (entity.CareerProfileTestSession, error) {
	if tx == nil {
		tx = r.db
	}

	var session entity.CareerProfileTestSession
	if err := tx.WithContext(ctx).Take(&session, "session_token = ?", token).Error; err != nil {
		return entity.CareerProfileTestSession{}, err
	}

	return session, nil
}

func (r *testSessionRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, timestamps map[string]interface{}) error {
	if tx == nil {
		tx = r.db
	}

	updates := map[string]interface{}{
		"status": status,
	}
	for k, v := range timestamps {
		updates[k] = v
	}

	return tx.WithContext(ctx).
		Model(&entity.CareerProfileTestSession{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *testSessionRepository) GetUserSessions(ctx context.Context, tx *gorm.DB, userID string) ([]entity.CareerProfileTestSession, error) {
	if tx == nil {
		tx = r.db
	}

	var sessions []entity.CareerProfileTestSession
	if err := tx.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("started_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}
