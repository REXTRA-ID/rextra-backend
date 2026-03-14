package service

import (
	"context"
	"math"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/modules/membership/repository"
)

type (
	SubscriptionCycleService interface {
		GetAllPaginated(ctx context.Context, request dto_request.SubscriptionCycleFilterRequest) (dto_response.PaginatedSubscriptionCycleResponse, error)
		GetByUserID(ctx context.Context, userId string) ([]dto_response.GetSubscriptionCycleResponse, error)
	}

	subscriptionCycleService struct {
		subscriptionRepo repository.SubscriptionCycleRepository
	}
)

func NewSubscriptionCycleService(subscriptionRepo repository.SubscriptionCycleRepository) SubscriptionCycleService {
	return &subscriptionCycleService{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *subscriptionCycleService) GetAllPaginated(ctx context.Context, request dto_request.SubscriptionCycleFilterRequest) (dto_response.PaginatedSubscriptionCycleResponse, error) {
	page := request.Page
	if page < 1 {
		page = 1
	}
	pageSize := request.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	cycles, err := s.subscriptionRepo.GetAllPaginated(ctx, nil, offset, pageSize)
	if err != nil {
		return dto_response.PaginatedSubscriptionCycleResponse{}, err
	}

	totalCount, err := s.subscriptionRepo.GetTotalCount(ctx, nil)
	if err != nil {
		return dto_response.PaginatedSubscriptionCycleResponse{}, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	var responses []dto_response.GetSubscriptionCycleResponse
	for _, cycle := range cycles {
		responses = append(responses, dto_response.GetSubscriptionCycleResponse{
			ID:             cycle.ID.String(),
			UserID:         cycle.UserID.String(),
			MembershipID:   cycle.MembershipID.String(),
			PlanName:       cycle.PlanName,
			DurationMonths: cycle.DurationMonths,
			AmountPaid:     cycle.AmountPaid,
			PaymentChannel: cycle.PaymentChannel,
			StartDate:      cycle.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:        cycle.EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	return dto_response.PaginatedSubscriptionCycleResponse{
		Data:       responses,
		Total:      int(totalCount),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *subscriptionCycleService) GetByUserID(ctx context.Context, userId string) ([]dto_response.GetSubscriptionCycleResponse, error) {
	cycles, err := s.subscriptionRepo.GetByUserID(ctx, nil, userId)
	if err != nil {
		return nil, err
	}

	var responses []dto_response.GetSubscriptionCycleResponse
	for _, cycle := range cycles {
		responses = append(responses, dto_response.GetSubscriptionCycleResponse{
			ID:             cycle.ID.String(),
			UserID:         cycle.UserID.String(),
			MembershipID:   cycle.MembershipID.String(),
			PlanName:       cycle.PlanName,
			DurationMonths: cycle.DurationMonths,
			AmountPaid:     cycle.AmountPaid,
			PaymentChannel: cycle.PaymentChannel,
			StartDate:      cycle.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:        cycle.EndDate.Format("2006-01-02 15:04:05"),
		})
	}

	return responses, nil
}
