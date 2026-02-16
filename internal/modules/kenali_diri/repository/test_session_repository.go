package repository

import (
	"context"
	"strings"
	"time"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	TestSessionFilters struct {
		TestGoal    string
		PersonaType string
		Status      string
		UserName    string
		StartDate   *time.Time
		EndDate     *time.Time
		SortBy      string
		Limit       int
		Offset      int
	}

	TestSessionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, session entity.CareerProfileTestSession) (entity.CareerProfileTestSession, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.CareerProfileTestSession, error)
		GetByToken(ctx context.Context, tx *gorm.DB, token string) (entity.CareerProfileTestSession, error)
		UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, timestamps map[string]interface{}) error
		GetUserSessions(ctx context.Context, tx *gorm.DB, userID string) ([]entity.CareerProfileTestSession, error)
		ListWithFilters(ctx context.Context, tx *gorm.DB, filters TestSessionFilters) ([]entity.CareerProfileTestSession, int64, error)
		BulkDelete(ctx context.Context, tx *gorm.DB, ids []int64) error
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
func (r *testSessionRepository) ListWithFilters(ctx context.Context, tx *gorm.DB, filters TestSessionFilters) ([]entity.CareerProfileTestSession, int64, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).Model(&entity.CareerProfileTestSession{})

	needJoinUser := filters.UserName != "" || strings.HasPrefix(filters.SortBy, "name")
	if needJoinUser {
		query = query.Joins("JOIN users ON users.id = careerprofile_test_sessions.user_id")
	}

	if filters.TestGoal != "" {
		query = query.Where("test_goal = ?", filters.TestGoal)
	}
	if filters.PersonaType != "" {
		query = query.Where("persona_type = ?", filters.PersonaType)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.UserName != "" {
		query = query.Where("LOWER(users.fullname) LIKE ?", "%"+strings.ToLower(filters.UserName)+"%")
	}
	if filters.StartDate != nil {
		query = query.Where("started_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("started_at <= ?", *filters.EndDate)
	}

	switch filters.SortBy {
	case "name_asc":
		query = query.Order("users.fullname ASC")
	case "name_desc":
		query = query.Order("users.fullname DESC")
	case "date_asc":
		query = query.Order("started_at ASC")
	case "date_desc":
		query = query.Order("started_at DESC")
	default:
		query = query.Order("started_at DESC")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var sessions []entity.CareerProfileTestSession
	if err := query.Preload("User").Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *testSessionRepository) BulkDelete(ctx context.Context, tx *gorm.DB, ids []int64) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&entity.CareerProfileTestSession{}).Error
}
