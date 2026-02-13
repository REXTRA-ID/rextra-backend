package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	dto_req "rextra-backend/internal/dto/request"
	dto_res "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/kenali_diri/repository"
)

type CareerProfileFeedbackService interface {
	GetStudentFeedbacks(ctx context.Context, req dto_req.GetStudentFeedbackRequest) (*dto_res.CareerProfileStudentFeedbackListResponse, error)
	GetExpertFeedbacks(ctx context.Context, req dto_req.GetExpertFeedbackRequest) (*dto_res.CareerProfileExpertFeedbackListResponse, error)
	GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error)
	GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error)
	GetStudentFeedbackStats(ctx context.Context, categoryID *int64, timeRange string) (dto_res.FeedbackStatsResponse, error)
}

type careerProfileFeedbackService struct {
	repo         repository.CareerProfileFeedbackRepository
	feedbackRepo repository.FeedbackRepository
}

func NewCareerProfileFeedbackService(repo repository.CareerProfileFeedbackRepository, feedbackRepo repository.FeedbackRepository) CareerProfileFeedbackService {
	return &careerProfileFeedbackService{
		repo:         repo,
		feedbackRepo: feedbackRepo,
	}
}

func (s *careerProfileFeedbackService) GetStudentFeedbacks(ctx context.Context, req dto_req.GetStudentFeedbackRequest) (*dto_res.CareerProfileStudentFeedbackListResponse, error) {
	items, total, err := s.repo.GetStudentFeedbackList(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int(total) / req.PageSize
		if int(total)%req.PageSize > 0 {
			totalPages++
		}
	}

	return &dto_res.CareerProfileStudentFeedbackListResponse{
		Items: items,
		Pagination: dto_res.CareerProfilePaginationMeta{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    req.Page < totalPages,
			HasPrev:    req.Page > 1,
		},
	}, nil
}

func (s *careerProfileFeedbackService) GetExpertFeedbacks(ctx context.Context, req dto_req.GetExpertFeedbackRequest) (*dto_res.CareerProfileExpertFeedbackListResponse, error) {
	items, total, err := s.repo.GetExpertFeedbackList(ctx, req)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int(total) / req.PageSize
		if int(total)%req.PageSize > 0 {
			totalPages++
		}
	}

	return &dto_res.CareerProfileExpertFeedbackListResponse{
		Items: items,
		Pagination: dto_res.CareerProfilePaginationMeta{
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
			HasNext:    req.Page < totalPages,
			HasPrev:    req.Page > 1,
		},
	}, nil
}

func (s *careerProfileFeedbackService) GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error) {
	return s.repo.GetExpertFeedbackDetail(ctx, id)
}

func (s *careerProfileFeedbackService) GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error) {
	return s.repo.GetFeedbackMetadata(ctx)
}

func (s *careerProfileFeedbackService) GetStudentFeedbackStats(ctx context.Context, categoryID *int64, timeRange string) (dto_res.FeedbackStatsResponse, error) {
	normalizedRange := normalizeTimeRange(timeRange)
	start := calculateStatsStart(normalizedRange)
	filters := repository.StudentFeedbackFilters{
		CategoryID: categoryID,
		StartDate:  start,
		TimeRange:  normalizedRange,
	}

	stats, err := s.feedbackRepo.GetStudentFeedbackStats(ctx, nil, filters)
	if err != nil {
		return dto_res.FeedbackStatsResponse{}, err
	}

	labels, testCounts, feedbackCounts := buildTrendSeries(stats.BucketGranularity, stats.TrendTests, stats.TrendFeedback)
	scoreDistribution := dto_res.ScoreDistributionSet{
		EaseOfUse:    convertScoreDistribution(stats.ScoreEase, stats.TotalFeedback),
		Relevance:    convertScoreDistribution(stats.ScoreRelevance, stats.TotalFeedback),
		Satisfaction: convertScoreDistribution(stats.ScoreSatisfaction, stats.TotalFeedback),
	}
	obstacleSummary := dto_res.ObstacleSummary{
		WithObstacles:    int(stats.WithObstacles),
		WithoutObstacles: int(stats.WithoutObstacles),
	}
	obstacleDistribution := convertObstacleDistribution(stats.ObstacleBreakdown, stats.TotalFeedback)
	sentimentComposition := convertSentimentComposition(stats.SentimentScores, stats.TotalFeedback)
	responseInfo := buildResponseInfo(stats.ResponseTotals)
	participationRate := responseInfo.ResponseRate

	return dto_res.FeedbackStatsResponse{
		TotalFeedback:     int(stats.TotalFeedback),
		AvgEaseOfUse:      stats.AvgEaseOfUse,
		AvgRelevance:      stats.AvgRelevance,
		AvgSatisfaction:   stats.AvgSatisfaction,
		ParticipationRate: participationRate,
		TrendData: dto_res.TrendChartData{
			Labels:         labels,
			TestCounts:     testCounts,
			FeedbackCounts: feedbackCounts,
		},
		ScoreDistribution:    scoreDistribution,
		ObstacleSummary:      obstacleSummary,
		ObstacleDistribution: obstacleDistribution,
		SentimentComposition: sentimentComposition,
		ResponseRateInfo:     responseInfo,
	}, nil
}

func normalizeTimeRange(value string) string {
	switch strings.ToLower(value) {
	case "bulanan":
		return "bulanan"
	case "sepanjang_waktu":
		return "sepanjang_waktu"
	default:
		return "mingguan"
	}
}

func calculateStatsStart(timeRange string) *time.Time {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	var start time.Time

	switch timeRange {
	case "bulanan":
		start = now.AddDate(0, 0, -27)
		return &start
	case "sepanjang_waktu":
		return nil
	default:
		start = now.AddDate(0, 0, -6)
		return &start
	}
}

func buildTrendSeries(bucket string, tests, feedback []repository.TimeBucketCount) ([]string, []int, []int) {
	bucketMap := make(map[int64]time.Time)
	for _, b := range tests {
		bucketMap[b.Bucket.Unix()] = b.Bucket
	}
	for _, b := range feedback {
		bucketMap[b.Bucket.Unix()] = b.Bucket
	}

	if len(bucketMap) == 0 {
		return []string{}, []int{}, []int{}
	}

	keys := make([]int64, 0, len(bucketMap))
	for k := range bucketMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	testMap := make(map[int64]int64)
	for _, b := range tests {
		testMap[b.Bucket.Unix()] = b.Count
	}

	feedbackMap := make(map[int64]int64)
	for _, b := range feedback {
		feedbackMap[b.Bucket.Unix()] = b.Count
	}

	labels := make([]string, 0, len(keys))
	testCounts := make([]int, 0, len(keys))
	feedbackCounts := make([]int, 0, len(keys))

	for _, key := range keys {
		labels = append(labels, formatTrendLabel(bucket, bucketMap[key]))
		testCounts = append(testCounts, int(testMap[key]))
		feedbackCounts = append(feedbackCounts, int(feedbackMap[key]))
	}

	return labels, testCounts, feedbackCounts
}

func formatTrendLabel(bucket string, t time.Time) string {
	switch bucket {
	case "week":
		start := t
		end := t.AddDate(0, 0, 6)
		return fmt.Sprintf("%s-%s", start.Format("02 Jan"), end.Format("02 Jan"))
	case "month":
		return t.Format("Jan 2006")
	default:
		return t.Format("Mon, 02 Jan")
	}
}

func convertScoreDistribution(entries []repository.ScoreCount, total int64) []dto_res.ScoreDistributionItem {
	result := make([]dto_res.ScoreDistributionItem, 0, 7)
	scoreMap := make(map[int]int64)
	for _, entry := range entries {
		scoreMap[entry.Score] = entry.Count
	}

	for score := 1; score <= 7; score++ {
		count := scoreMap[score]
		result = append(result, dto_res.ScoreDistributionItem{
			Score:      score,
			Count:      int(count),
			Percentage: percentage(count, total),
		})
	}

	return result
}

func convertObstacleDistribution(entries []repository.ObstacleCount, total int64) []dto_res.ObstacleDistributionItem {
	result := make([]dto_res.ObstacleDistributionItem, 0, len(entries))
	for _, entry := range entries {
		result = append(result, dto_res.ObstacleDistributionItem{
			Name:       entry.Name,
			Count:      int(entry.Count),
			Percentage: percentage(entry.Count, total),
		})
	}

	return result
}

func convertSentimentComposition(data map[string]repository.SentimentCount, total int64) []dto_res.SentimentCompositionItem {
	metrics := []struct {
		Key   string
		Label string
	}{
		{Key: "ease_of_use", Label: "Kemudahan Tes"},
		{Key: "relevance", Label: "Relevansi Rekomendasi"},
		{Key: "satisfaction", Label: "Kepuasan Fitur"},
	}

	result := make([]dto_res.SentimentCompositionItem, 0, len(metrics))
	for _, metric := range metrics {
		count := data[metric.Key]
		result = append(result, dto_res.SentimentCompositionItem{
			Metric:   metric.Label,
			Negative: percentage(count.Negative, total),
			Neutral:  percentage(count.Neutral, total),
			Positive: percentage(count.Positive, total),
		})
	}

	return result
}

func buildResponseInfo(totals repository.ResponseRateTotals) dto_res.ResponseRateInfo {
	notFilled := totals.CompletedTests - totals.FeedbackCount
	if notFilled < 0 {
		notFilled = 0
	}

	rate := 0.0
	if totals.CompletedTests > 0 {
		rate = math.Round((float64(totals.FeedbackCount)/float64(totals.CompletedTests))*100*100) / 100
	}

	return dto_res.ResponseRateInfo{
		FeedbackCount:  int(totals.FeedbackCount),
		CompletedTests: int(totals.CompletedTests),
		NotFilled:      int(notFilled),
		ResponseRate:   rate,
	}
}

func percentage(part, total int64) float64 {
	if total <= 0 {
		return 0
	}

	return math.Round((float64(part)/float64(total))*100*100) / 100
}
