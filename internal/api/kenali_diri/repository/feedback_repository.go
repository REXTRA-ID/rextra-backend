package repository

import (
	"context"
	"encoding/json"
	"strings"

	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	StudentFeedbackFilters struct {
		CategoryID   *int64
		UserName     string
		HasObstacles *bool
		SortBy       string
		Limit        int
		Offset       int
	}

	ExpertFeedbackFilters struct {
		CategoryID *int64
		ExpertName string
		TopNStatus string
		SortBy     string
		Limit      int
		Offset     int
	}

	FeedbackRepository interface {
		ListStudentFeedback(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) ([]entity.StudentFeedback, int64, error)
		GetStudentFeedbackByID(ctx context.Context, tx *gorm.DB, id int64) (entity.StudentFeedback, error)
		GetStudentFeedbackStats(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (map[string]interface{}, error)
		ListExpertFeedback(ctx context.Context, tx *gorm.DB, filters ExpertFeedbackFilters) ([]entity.ExpertFeedback, int64, error)
		GetExpertFeedbackByID(ctx context.Context, tx *gorm.DB, id int64) (entity.ExpertFeedback, error)
	}

	feedbackRepository struct {
		db *gorm.DB
	}
)

func NewFeedback(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) ListStudentFeedback(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) ([]entity.StudentFeedback, int64, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).Model(&entity.StudentFeedback{}).Preload("User").Preload("TestCategory")

	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}
	if filters.UserName != "" {
		query = query.Joins("JOIN users ON users.id = student_feedback.user_id").
			Where("LOWER(users.fullname) LIKE ?", "%"+strings.ToLower(filters.UserName)+"%")
	}
	if filters.HasObstacles != nil {
		if *filters.HasObstacles {
			query = query.Where("jsonb_array_length(obstacles) > 0")
		} else {
			query = query.Where("jsonb_array_length(obstacles) = 0")
		}
	}

	switch filters.SortBy {
	case "name_asc":
		query = query.Order("users.fullname ASC")
	case "name_desc":
		query = query.Order("users.fullname DESC")
	case "date_asc":
		query = query.Order("submitted_at ASC")
	case "date_desc":
		query = query.Order("submitted_at DESC")
	default:
		query = query.Order("submitted_at DESC")
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

	var feedbacks []entity.StudentFeedback
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

func (r *feedbackRepository) GetStudentFeedbackByID(ctx context.Context, tx *gorm.DB, id int64) (entity.StudentFeedback, error) {
	if tx == nil {
		tx = r.db
	}

	var feedback entity.StudentFeedback
	if err := tx.WithContext(ctx).Preload("User").Preload("TestCategory").Take(&feedback, "id = ?", id).Error; err != nil {
		return entity.StudentFeedback{}, err
	}

	return feedback, nil
}

func (r *feedbackRepository) GetStudentFeedbackStats(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (map[string]interface{}, error) {
	if tx == nil {
		tx = r.db
	}

	type stats struct {
		Total            int64
		AvgEase          float64
		AvgRelevance     float64
		AvgSatisfaction  float64
		WithObstacles    int64
		WithoutObstacles int64
	}

	query := tx.WithContext(ctx).Model(&entity.StudentFeedback{})
	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}

	var s stats
	if err := query.Select(
		"COUNT(*) as total",
		"AVG(ease_of_use_score) as avg_ease",
		"AVG(relevance_score) as avg_relevance",
		"AVG(satisfaction_score) as avg_satisfaction",
		"SUM(CASE WHEN jsonb_array_length(obstacles) > 0 THEN 1 ELSE 0 END) as with_obstacles",
		"SUM(CASE WHEN jsonb_array_length(obstacles) = 0 THEN 1 ELSE 0 END) as without_obstacles",
	).Scan(&s).Error; err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"total_feedback":     s.Total,
		"avg_ease_of_use":    s.AvgEase,
		"avg_relevance":      s.AvgRelevance,
		"avg_satisfaction":   s.AvgSatisfaction,
		"participation_rate": float64(0),
		"trend_data": map[string]interface{}{
			"dates":            []string{},
			"test_counts":      []int{},
			"feedback_counts":  []int{},
			"obstacle_present": s.WithObstacles,
			"obstacle_empty":   s.WithoutObstacles,
		},
	}

	return response, nil
}

func (r *feedbackRepository) ListExpertFeedback(ctx context.Context, tx *gorm.DB, filters ExpertFeedbackFilters) ([]entity.ExpertFeedback, int64, error) {
	if tx == nil {
		tx = r.db
	}

	query := tx.WithContext(ctx).Model(&entity.ExpertFeedback{}).Preload("TestCategory")

	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}
	if filters.ExpertName != "" {
		query = query.Where("LOWER(expert_name) LIKE ?", "%"+strings.ToLower(filters.ExpertName)+"%")
	}

	switch filters.SortBy {
	case "name_asc":
		query = query.Order("expert_name ASC")
	case "name_desc":
		query = query.Order("expert_name DESC")
	case "date_asc":
		query = query.Order("submitted_at ASC")
	case "date_desc":
		query = query.Order("submitted_at DESC")
	default:
		query = query.Order("submitted_at DESC")
	}

	var base []entity.ExpertFeedback
	if err := query.Find(&base).Error; err != nil {
		return nil, 0, err
	}

	filtered := make([]entity.ExpertFeedback, 0, len(base))
	for _, fb := range base {
		status := deriveTopNStatus(fb)
		if filters.TopNStatus == "" || strings.EqualFold(status, filters.TopNStatus) {
			filtered = append(filtered, fb)
		}
	}

	total := int64(len(filtered))
	start := filters.Offset
	if start < 0 {
		start = 0
	}
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + filters.Limit
	if filters.Limit == 0 || end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (r *feedbackRepository) GetExpertFeedbackByID(ctx context.Context, tx *gorm.DB, id int64) (entity.ExpertFeedback, error) {
	if tx == nil {
		tx = r.db
	}

	var feedback entity.ExpertFeedback
	if err := tx.WithContext(ctx).
		Preload("TestSession").
		Preload("TestCategory").
		Take(&feedback, "id = ?", id).Error; err != nil {
		return entity.ExpertFeedback{}, err
	}

	return feedback, nil
}

func deriveTopNStatus(fb entity.ExpertFeedback) string {
	if len(fb.TopFiveProfessions) == 0 {
		return "not_found"
	}

	var topFive []string
	_ = json.Unmarshal(fb.TopFiveProfessions, &topFive)

	prof := strings.TrimSpace(strings.ToLower(fb.Profession))
	for idx, item := range topFive {
		item = strings.TrimSpace(strings.ToLower(item))
		if item == "" || prof == "" {
			continue
		}
		if item == prof {
			switch idx {
			case 0:
				return "P1"
			case 1:
				return "P2"
			case 2, 3, 4:
				return "P3-5"
			default:
				return "P3-5"
			}
		}
	}

	return "not_found"
}
