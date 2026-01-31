package repository

import (
	"context"
	"strings"
	"time"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	KenalidiriHistoryFilters struct {
		CategoryID *int64
		Status     string
		UserName   string
		StartDate  *time.Time
		EndDate    *time.Time
		SortBy     string
		Limit      int
		Offset     int
	}

	KenalidiriHistoryRepository interface {
		Create(ctx context.Context, tx *gorm.DB, history entity.KenaliDiriHistory) (entity.KenaliDiriHistory, error)
		GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.KenaliDiriHistory, error)
		UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, completedAt *time.Time) error
		ListWithFilters(ctx context.Context, tx *gorm.DB, filters KenalidiriHistoryFilters) ([]entity.KenaliDiriHistory, int64, error)
		GetUserHistory(ctx context.Context, tx *gorm.DB, userID string, categoryID *int64) ([]entity.KenaliDiriHistory, error)
		BulkDelete(ctx context.Context, tx *gorm.DB, ids []int64) error
		CountByStatus(ctx context.Context, tx *gorm.DB, categoryID *int64) (map[string]int64, error)
	}

	kenalidiriHistoryRepository struct {
		db *gorm.DB
	}
)

func NewKenalidiriHistory(db *gorm.DB) KenalidiriHistoryRepository {
	return &kenalidiriHistoryRepository{db: db}
}

func (r *kenalidiriHistoryRepository) Create(ctx context.Context, tx *gorm.DB, history entity.KenaliDiriHistory) (entity.KenaliDiriHistory, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&history).Error; err != nil {
		return history, err
	}

	return history, nil
}

func (r *kenalidiriHistoryRepository) GetByID(ctx context.Context, tx *gorm.DB, id int64) (entity.KenaliDiriHistory, error) {
	if tx == nil {
		tx = r.db
	}

	var history entity.KenaliDiriHistory
	err := tx.WithContext(ctx).
		Preload("User").
		Preload("TestCategory").
		Take(&history, "id = ?", id).Error
	if err != nil {
		return entity.KenaliDiriHistory{}, err
	}

	return history, nil
}

func (r *kenalidiriHistoryRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id int64, status string, completedAt *time.Time) error {
	if tx == nil {
		tx = r.db
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if completedAt != nil {
		updates["completed_at"] = completedAt
	}

	return tx.WithContext(ctx).
		Model(&entity.KenaliDiriHistory{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *kenalidiriHistoryRepository) ListWithFilters(ctx context.Context, tx *gorm.DB, filters KenalidiriHistoryFilters) ([]entity.KenaliDiriHistory, int64, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).Model(&entity.KenaliDiriHistory{})

	needJoinUser := filters.UserName != "" || strings.HasPrefix(filters.SortBy, "name")
	if needJoinUser {
		query = query.Joins("JOIN users ON users.id = kenalidiri_history.user_id")
	}

	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
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

	var histories []entity.KenaliDiriHistory
	if err := query.Preload("User").Preload("TestCategory").Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *kenalidiriHistoryRepository) GetUserHistory(ctx context.Context, tx *gorm.DB, userID string, categoryID *int64) ([]entity.KenaliDiriHistory, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).
		Preload("TestCategory").
		Where("user_id = ?", userID)
	if categoryID != nil {
		query = query.Where("test_category_id = ?", *categoryID)
	}

	var histories []entity.KenaliDiriHistory
	if err := query.Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *kenalidiriHistoryRepository) BulkDelete(ctx context.Context, tx *gorm.DB, ids []int64) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&entity.KenaliDiriHistory{}).Error
}

func (r *kenalidiriHistoryRepository) CountByStatus(ctx context.Context, tx *gorm.DB, categoryID *int64) (map[string]int64, error) {
	if tx == nil {
		tx = r.db
	}

	type result struct {
		Status string
		Total  int64
	}

	var results []result
	query := tx.WithContext(ctx).Model(&entity.KenaliDiriHistory{})
	if categoryID != nil {
		query = query.Where("test_category_id = ?", *categoryID)
	}

	if err := query.Select("status, COUNT(*) as total").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}

	output := make(map[string]int64)
	for _, r := range results {
		output[r.Status] = r.Total
	}

	return output, nil
}
