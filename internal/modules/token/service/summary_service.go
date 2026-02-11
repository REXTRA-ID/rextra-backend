package service

import (
	"context"
	"time"

	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/token/repository"
	myerror "rextra-backend/internal/pkg/error"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const (
	dateLayout        = "2006-01-02"
	hoursPerDay       = 24
	fullGrowthPercent = 100.0
)

type (
	TokenSummaryService interface {
		GetKPISummary(ctx context.Context, startDate, endDate string) (dto_response.TokenKPIOverviewDTOResponse, error)
		GetGraphTrendByDirection(ctx context.Context, startDate, endDate string) ([]dto_response.TokenTrendSummaryDTOResponse, error)
		GetGraphTrendBySourceType(ctx context.Context, sourceType entity.TokenSourceType, startDate, endDate string) ([]dto_response.TokenSourceTrendDTOResponse, error)
	}

	tokenSummaryService struct {
		tokenLedgerRepo repository.TokenLedgerRepository
		topupRepo       repository.TopupTransactionRepository
		db              *gorm.DB
	}

	// kpiMetrics holds all KPI metric values
	kpiMetrics struct {
		tokenIn              int
		tokenOut             int
		tokenUsage           int
		membershipAllocation int
		topupSuccess         int
	}

	// periodMetrics holds metrics for current and previous periods
	periodMetrics struct {
		current  kpiMetrics
		previous kpiMetrics
	}
)

func NewTokenSummaryService(tokenLedgerRepo repository.TokenLedgerRepository, topupRepo repository.TopupTransactionRepository, db *gorm.DB) TokenSummaryService {
	return &tokenSummaryService{
		tokenLedgerRepo: tokenLedgerRepo,
		topupRepo:       topupRepo,
		db:              db,
	}
}

func (s *tokenSummaryService) GetKPISummary(ctx context.Context, startDate, endDate string) (dto_response.TokenKPIOverviewDTOResponse, error) {
	prevStartDate, prevEndDate := calculatePreviousPeriod(startDate, endDate)

	metrics, err := s.fetchAllMetrics(ctx, startDate, endDate, prevStartDate, prevEndDate)
	if err != nil {
		return dto_response.TokenKPIOverviewDTOResponse{}, err
	}

	return s.buildKPIResponse(metrics), nil
}
func (s *tokenSummaryService) GetGraphTrendByDirection(ctx context.Context, startDate, endDate string) ([]dto_response.TokenTrendSummaryDTOResponse, error) {
	results, err := s.tokenLedgerRepo.TrendByDirection(ctx, nil, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summaryMap := make(map[string]*dto_response.TokenTrendSummaryDTOResponse)
	var dates []string

	for _, res := range results {
		dateStr := res.Date.Format("2006-01-02")
		summary, exists := summaryMap[dateStr]
		if !exists {
			summary = &dto_response.TokenTrendSummaryDTOResponse{Date: dateStr}
			summaryMap[dateStr] = summary
			dates = append(dates, dateStr)
		}

		if res.Direction == entity.DirectionIN {
			summary.TokenIn += res.Amount
		} else if res.Direction == entity.DirectionOUT {
			summary.TokenOut += res.Amount
		}
	}

	var finalSummary []dto_response.TokenTrendSummaryDTOResponse
	for _, dateStr := range dates {
		summary := summaryMap[dateStr]
		summary.Net = summary.TokenIn - summary.TokenOut
		finalSummary = append(finalSummary, *summary)
	}

	return finalSummary, nil
}

func (s *tokenSummaryService) GetGraphTrendBySourceType(ctx context.Context, sourceType entity.TokenSourceType, startDate, endDate string) ([]dto_response.TokenSourceTrendDTOResponse, error) {
	if !sourceType.IsValid() {
		return nil, myerror.New("Invalid Source type", myerror.Error_InvalidRequest)
	}

	results, err := s.tokenLedgerRepo.TrendBySourceType(ctx, nil, string(sourceType), startDate, endDate)
	if err != nil {
		return nil, err
	}

	var graph []dto_response.TokenSourceTrendDTOResponse
	for _, res := range results {
		dateStr := res.Date.Format("2006-01-02")
		graph = append(graph, dto_response.TokenSourceTrendDTOResponse{
			Date:  dateStr,
			Value: res.Amount,
		})
	}

	return graph, nil
}

// ===== HELPER FUNCTIONS ======== //

// fetchAllMetrics retrieves all KPI metrics concurrently for both periods
func (s *tokenSummaryService) fetchAllMetrics(ctx context.Context, startDate, endDate, prevStartDate, prevEndDate string) (periodMetrics, error) {
	var metrics periodMetrics
	g, ctx := errgroup.WithContext(ctx)

	s.fetchMetricPair(g, ctx, &metrics.current.tokenIn, &metrics.previous.tokenIn, s.countTokensByDirection("IN"), startDate, endDate, prevStartDate, prevEndDate)
	s.fetchMetricPair(g, ctx, &metrics.current.tokenOut, &metrics.previous.tokenOut, s.countTokensByDirection("OUT"), startDate, endDate, prevStartDate, prevEndDate)
	s.fetchMetricPair(g, ctx, &metrics.current.tokenUsage, &metrics.previous.tokenUsage, s.countTokensByType("USAGE"), startDate, endDate, prevStartDate, prevEndDate)
	s.fetchMetricPair(g, ctx, &metrics.current.membershipAllocation, &metrics.previous.membershipAllocation, s.countTokensByType("MEMBERSHIP"), startDate, endDate, prevStartDate, prevEndDate)
	s.fetchMetricPair(g, ctx, &metrics.current.topupSuccess, &metrics.previous.topupSuccess, s.countTopupsByStatus("SUCCESS"), startDate, endDate, prevStartDate, prevEndDate)

	if err := g.Wait(); err != nil {
		return periodMetrics{}, err
	}

	return metrics, nil
}

// fetchMetricPair fetches current and previous values for a metric concurrently
func (s *tokenSummaryService) fetchMetricPair(g *errgroup.Group, ctx context.Context, currentValue, previousValue *int, fetcher func(ctx context.Context, start, end string) (int, error), startDate, endDate, prevStartDate, prevEndDate string) {
	g.Go(func() error {
		val, err := fetcher(ctx, startDate, endDate)
		if err != nil {
			return err
		}
		*currentValue = val
		return nil
	})

	g.Go(func() error {
		val, err := fetcher(ctx, prevStartDate, prevEndDate)
		if err != nil {
			return err
		}
		*previousValue = val
		return nil
	})
}

// countTokensByDirection returns a function that counts tokens by direction
func (s *tokenSummaryService) countTokensByDirection(direction string) func(context.Context, string, string) (int, error) {
	return func(ctx context.Context, start, end string) (int, error) {
		return s.tokenLedgerRepo.CountTokenByDirection(ctx, nil, direction, start, end)
	}
}

// countTokensByType returns a function that counts tokens by type
func (s *tokenSummaryService) countTokensByType(tokenType string) func(context.Context, string, string) (int, error) {
	return func(ctx context.Context, start, end string) (int, error) {
		return s.tokenLedgerRepo.CountTokenByType(ctx, nil, tokenType, start, end)
	}
}

// countTopupsByStatus returns a function that counts topups by status
func (s *tokenSummaryService) countTopupsByStatus(status string) func(context.Context, string, string) (int, error) {
	return func(ctx context.Context, start, end string) (int, error) {
		return s.topupRepo.CountByStatus(ctx, nil, entity.TopupStatus(status), start, end)
	}
}

// buildKPIResponse constructs the response DTO from metrics
func (s *tokenSummaryService) buildKPIResponse(metrics periodMetrics) dto_response.TokenKPIOverviewDTOResponse {
	curr := metrics.current
	prev := metrics.previous

	netFlowCurr := curr.tokenIn - curr.tokenOut
	netFlowPrev := prev.tokenIn - prev.tokenOut

	return dto_response.TokenKPIOverviewDTOResponse{
		TokenIn:              s.createKPI(int64(curr.tokenIn), "Token Masuk", curr.tokenIn, prev.tokenIn),
		TokenOut:             s.createKPI(int64(curr.tokenOut), "Token Keluar", curr.tokenOut, prev.tokenOut),
		Netflow:              s.createKPI(int64(netFlowCurr), "Net Flow", netFlowCurr, netFlowPrev),
		TokenUsage:           s.createKPI(int64(curr.tokenUsage), "Pemakaian Token", curr.tokenUsage, prev.tokenUsage),
		MembershipAllocation: s.createKPI(int64(curr.membershipAllocation), "Alokasi Membership", curr.membershipAllocation, prev.membershipAllocation),
		TopupSuccess:         s.createKPI(int64(curr.topupSuccess), "Topup Berhasil", curr.topupSuccess, prev.topupSuccess),
	}
}

// createKPI creates a KPI object with growth calculation
func (s *tokenSummaryService) createKPI(value int64, label string, current, previous int) dto_response.KPI {
	return dto_response.KPI{
		Value:      value,
		Label:      label,
		Percentage: calculateGrowth(current, previous),
	}
}

// calculatePreviousPeriod calculates the previous period based on current period duration
func calculatePreviousPeriod(start, end string) (string, string) {
	startDate, err := time.Parse(dateLayout, start)
	if err != nil {
		return "", ""
	}

	endDate, err := time.Parse(dateLayout, end)
	if err != nil {
		return "", ""
	}

	periodDuration := endDate.Sub(startDate)
	previousEnd := startDate.Add(-hoursPerDay * time.Hour)
	previousStart := previousEnd.Add(-periodDuration)

	return previousStart.Format(dateLayout), previousEnd.Format(dateLayout)
}

// calculateGrowth calculates percentage growth between current and previous values
func calculateGrowth(current, previous int) float64 {
	if previous == 0 {
		if current > 0 {
			return fullGrowthPercent
		}
		return 0.0
	}
	return ((float64(current) - float64(previous)) / float64(previous)) * 100
}
