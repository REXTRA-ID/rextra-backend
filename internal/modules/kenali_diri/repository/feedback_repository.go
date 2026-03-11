package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
		StartDate    *time.Time
		EndDate      *time.Time
		TimeRange    string
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
		GetStudentFeedbackStats(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (StudentFeedbackStatsResult, error)
		ListExpertFeedback(ctx context.Context, tx *gorm.DB, filters ExpertFeedbackFilters) ([]entity.ExpertFeedback, int64, error)
		GetExpertFeedbackByID(ctx context.Context, tx *gorm.DB, id int64) (entity.ExpertFeedback, error)
	}

	feedbackRepository struct {
		db *gorm.DB
	}

	TimeBucketCount struct {
		Bucket time.Time `json:"bucket"`
		Count  int64     `json:"count"`
	}

	ScoreCount struct {
		Score int   `json:"score"`
		Count int64 `json:"count"`
	}

	ObstacleCount struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}

	SentimentCount struct {
		Negative int64 `json:"negative"`
		Neutral  int64 `json:"neutral"`
		Positive int64 `json:"positive"`
	}

	ResponseRateTotals struct {
		FeedbackCount  int64
		CompletedTests int64
	}

	StudentFeedbackStatsResult struct {
		TotalFeedback     int64
		AvgEaseOfUse      float64
		AvgRelevance      float64
		AvgSatisfaction   float64
		WithObstacles     int64
		WithoutObstacles  int64
		TrendFeedback     []TimeBucketCount
		TrendTests        []TimeBucketCount
		ScoreEase         []ScoreCount
		ScoreRelevance    []ScoreCount
		ScoreSatisfaction []ScoreCount
		ObstacleBreakdown []ObstacleCount
		SentimentScores   map[string]SentimentCount
		ResponseTotals    ResponseRateTotals
		BucketGranularity string
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

func (r *feedbackRepository) GetStudentFeedbackStats(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (StudentFeedbackStatsResult, error) {
	if tx == nil {
		tx = r.db
	}

	result := StudentFeedbackStatsResult{
		SentimentScores: make(map[string]SentimentCount),
	}

	bucket := bucketFromRange(filters.TimeRange)
	result.BucketGranularity = bucket

	baseQuery := r.applyStudentFeedbackFilters(tx.WithContext(ctx).Model(&entity.StudentFeedback{}), filters)
	var overview struct {
		Total            int64
		AvgEase          float64
		AvgRelevance     float64
		AvgSatisfaction  float64
		WithObstacles    int64
		WithoutObstacles int64
	}

	if err := baseQuery.Select(
		"COUNT(*) as total",
		"COALESCE(AVG(ease_of_use_score), 0) as avg_ease",
		"COALESCE(AVG(relevance_score), 0) as avg_relevance",
		"COALESCE(AVG(satisfaction_score), 0) as avg_satisfaction",
		"SUM(CASE WHEN jsonb_array_length(obstacles) > 0 THEN 1 ELSE 0 END) as with_obstacles",
		"SUM(CASE WHEN jsonb_array_length(obstacles) = 0 THEN 1 ELSE 0 END) as without_obstacles",
	).Scan(&overview).Error; err != nil {
		return result, err
	}

	result.TotalFeedback = overview.Total
	result.AvgEaseOfUse = overview.AvgEase
	result.AvgRelevance = overview.AvgRelevance
	result.AvgSatisfaction = overview.AvgSatisfaction
	result.WithObstacles = overview.WithObstacles
	result.WithoutObstacles = overview.WithoutObstacles
	result.ResponseTotals.FeedbackCount = overview.Total

	var err error
	if result.ScoreEase, err = r.fetchScoreDistribution(ctx, tx, filters, "ease_of_use_score"); err != nil {
		return result, err
	}
	if result.ScoreRelevance, err = r.fetchScoreDistribution(ctx, tx, filters, "relevance_score"); err != nil {
		return result, err
	}
	if result.ScoreSatisfaction, err = r.fetchScoreDistribution(ctx, tx, filters, "satisfaction_score"); err != nil {
		return result, err
	}

	if result.TrendFeedback, err = r.fetchFeedbackTrend(ctx, tx, filters, bucket); err != nil {
		return result, err
	}
	if result.TrendTests, err = r.fetchHistoryTrend(ctx, tx, filters, bucket); err != nil {
		return result, err
	}

	if result.ResponseTotals.CompletedTests, err = r.countCompletedTests(ctx, tx, filters); err != nil {
		return result, err
	}

	if result.ObstacleBreakdown, err = r.fetchObstacleBreakdown(ctx, tx, filters); err != nil {
		return result, err
	}

	if result.SentimentScores, err = r.fetchSentimentCounts(ctx, tx, filters); err != nil {
		return result, err
	}

	return result, nil
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

func (r *feedbackRepository) applyStudentFeedbackFilters(query *gorm.DB, filters StudentFeedbackFilters) *gorm.DB {
	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}
	if filters.StartDate != nil {
		query = query.Where("submitted_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("submitted_at <= ?", *filters.EndDate)
	}
	return query
}

func bucketFromRange(timeRange string) string {
	switch strings.ToLower(timeRange) {
	case "bulanan":
		return "week"
	case "sepanjang_waktu":
		return "month"
	default:
		return "day"
	}
}

func (r *feedbackRepository) fetchScoreDistribution(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters, column string) ([]ScoreCount, error) {
	query := r.applyStudentFeedbackFilters(tx.WithContext(ctx).Model(&entity.StudentFeedback{}), filters)
	selectExpr := fmt.Sprintf("%s AS score, COUNT(*) AS count", column)

	var rows []ScoreCount
	if err := query.Select(selectExpr).Group(column).Order("score").Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *feedbackRepository) fetchFeedbackTrend(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters, bucket string) ([]TimeBucketCount, error) {
	query := r.applyStudentFeedbackFilters(tx.WithContext(ctx).Model(&entity.StudentFeedback{}), filters)
	expr := fmt.Sprintf("date_trunc('%s', submitted_at)", bucket)

	var rows []TimeBucketCount
	if err := query.Select(expr + " AS bucket, COUNT(*) AS count").Group("bucket").Order("bucket").Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *feedbackRepository) fetchHistoryTrend(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters, bucket string) ([]TimeBucketCount, error) {
	query := tx.WithContext(ctx).Model(&entity.KenaliDiriHistory{})
	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}
	if filters.StartDate != nil {
		query = query.Where("started_at >= ?", *filters.StartDate)
	}

	expr := fmt.Sprintf("date_trunc('%s', started_at)", bucket)
	var rows []TimeBucketCount
	if err := query.Select(expr + " AS bucket, COUNT(*) AS count").Group("bucket").Order("bucket").Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *feedbackRepository) countCompletedTests(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (int64, error) {
	query := tx.WithContext(ctx).Model(&entity.KenaliDiriHistory{}).Where("status = ?", "completed").Where("completed_at IS NOT NULL")
	if filters.CategoryID != nil {
		query = query.Where("test_category_id = ?", *filters.CategoryID)
	}
	if filters.StartDate != nil {
		query = query.Where("completed_at >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("completed_at <= ?", *filters.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *feedbackRepository) fetchObstacleBreakdown(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) ([]ObstacleCount, error) {
	sql := "SELECT TRIM(value) AS name, COUNT(*) AS count FROM student_feedback, jsonb_array_elements_text(obstacles) AS value WHERE jsonb_array_length(obstacles) > 0 AND LOWER(TRIM(value)) <> 'tidak ada kendala' AND TRIM(value) <> ''"
	var args []interface{}

	if filters.CategoryID != nil {
		sql += " AND test_category_id = ?"
		args = append(args, *filters.CategoryID)
	}
	if filters.StartDate != nil {
		sql += " AND submitted_at >= ?"
		args = append(args, *filters.StartDate)
	}
	if filters.EndDate != nil {
		sql += " AND submitted_at <= ?"
		args = append(args, *filters.EndDate)
	}

	sql += " GROUP BY name ORDER BY count DESC"

	var rows []ObstacleCount
	if err := tx.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *feedbackRepository) fetchSentimentCounts(ctx context.Context, tx *gorm.DB, filters StudentFeedbackFilters) (map[string]SentimentCount, error) {
	query := r.applyStudentFeedbackFilters(tx.WithContext(ctx).Model(&entity.StudentFeedback{}), filters)

	var row struct {
		EaseNeg int64
		EaseNeu int64
		EasePos int64
		RelNeg  int64
		RelNeu  int64
		RelPos  int64
		SatNeg  int64
		SatNeu  int64
		SatPos  int64
	}

	if err := query.Select(`
		SUM(CASE WHEN ease_of_use_score BETWEEN 1 AND 3 THEN 1 ELSE 0 END) AS ease_neg,
		SUM(CASE WHEN ease_of_use_score = 4 THEN 1 ELSE 0 END) AS ease_neu,
		SUM(CASE WHEN ease_of_use_score BETWEEN 5 AND 7 THEN 1 ELSE 0 END) AS ease_pos,
		SUM(CASE WHEN relevance_score BETWEEN 1 AND 3 THEN 1 ELSE 0 END) AS rel_neg,
		SUM(CASE WHEN relevance_score = 4 THEN 1 ELSE 0 END) AS rel_neu,
		SUM(CASE WHEN relevance_score BETWEEN 5 AND 7 THEN 1 ELSE 0 END) AS rel_pos,
		SUM(CASE WHEN satisfaction_score BETWEEN 1 AND 3 THEN 1 ELSE 0 END) AS sat_neg,
		SUM(CASE WHEN satisfaction_score = 4 THEN 1 ELSE 0 END) AS sat_neu,
		SUM(CASE WHEN satisfaction_score BETWEEN 5 AND 7 THEN 1 ELSE 0 END) AS sat_pos
	`).Scan(&row).Error; err != nil {
		return nil, err
	}

	return map[string]SentimentCount{
		"ease_of_use": {Negative: row.EaseNeg, Neutral: row.EaseNeu, Positive: row.EasePos},
		"relevance":   {Negative: row.RelNeg, Neutral: row.RelNeu, Positive: row.RelPos},
		"satisfaction": {
			Negative: row.SatNeg,
			Neutral:  row.SatNeu,
			Positive: row.SatPos,
		},
	}, nil
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
