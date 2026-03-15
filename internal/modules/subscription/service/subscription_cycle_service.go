package service

import (
	"context"
	"errors"
	"math"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/subscription/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	SubscriptionCycleService interface {
		GetAll(ctx context.Context, req dto_request.SubscriptionCycleFilterRequest) (dto_response.PaginatedSubscriptionCycleResponse, error)
		GetById(ctx context.Context, cycleID string) (dto_response.GetSubscriptionCycleResponse, error)
		GetByUserID(ctx context.Context, userID string) ([]dto_response.GetSubscriptionCycleResponse, error)
		GetMyCycles(ctx context.Context, userID string) ([]dto_response.GetSubscriptionCycleResponse, error)
	}

	subscriptionCycleService struct {
		cycleRepo      repository.SubscriptionCycleRepository
		membershipRepo membershipRepo.UserMembershipRepository
	}
)

func NewSubscriptionCycleService(cycleRepo repository.SubscriptionCycleRepository, membershipRepo membershipRepo.UserMembershipRepository) SubscriptionCycleService {
	return &subscriptionCycleService{cycleRepo: cycleRepo, membershipRepo: membershipRepo}
}

func (s *subscriptionCycleService) GetAll(ctx context.Context, req dto_request.SubscriptionCycleFilterRequest) (dto_response.PaginatedSubscriptionCycleResponse, error) {
	page, pageSize := req.Page, req.PageSize
	if pageSize <= 0 { pageSize = 20 }
	offset := page * pageSize
	var dateFrom, dateTo time.Time
	if req.DateFrom != "" { dateFrom, _ = time.Parse("2006-01-02", req.DateFrom) }
	if req.DateTo != "" { dateTo, _ = time.Parse("2006-01-02", req.DateTo) }
	var filterUserID uuid.UUID
	if req.UserID != "" { filterUserID, _ = uuid.Parse(req.UserID) }

	cycles, total, err := s.cycleRepo.GetAllPaginated(ctx, nil, offset, pageSize, req.PlanName, req.Status, dateFrom, dateTo, filterUserID)
	if err != nil { return dto_response.PaginatedSubscriptionCycleResponse{}, myerror.DatabaseError(err) }
	var data []dto_response.GetSubscriptionCycleResponse
	for _, c := range cycles { data = append(data, toCycleResponse(c)) }
	return dto_response.PaginatedSubscriptionCycleResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *subscriptionCycleService) GetById(ctx context.Context, cycleID string) (dto_response.GetSubscriptionCycleResponse, error) {
	parsedID, err := uuid.Parse(cycleID)
	if err != nil { return dto_response.GetSubscriptionCycleResponse{}, myerror.InvalidRequest(err) }
	cycle, err := s.cycleRepo.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetSubscriptionCycleResponse{}, myerror.RecordNotFound("subscription cycle") }
		return dto_response.GetSubscriptionCycleResponse{}, myerror.DatabaseError(err)
	}
	return toCycleResponse(cycle), nil
}

func (s *subscriptionCycleService) GetByUserID(ctx context.Context, userID string) ([]dto_response.GetSubscriptionCycleResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return nil, myerror.InvalidRequest(err) }
	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	if err != nil { return nil, myerror.RecordNotFound("membership") }
	cycles, err := s.cycleRepo.GetAllByMembershipID(ctx, nil, membership.ID)
	if err != nil { return nil, myerror.DatabaseError(err) }
	var res []dto_response.GetSubscriptionCycleResponse
	for _, c := range cycles { res = append(res, toCycleResponse(c)) }
	return res, nil
}

func (s *subscriptionCycleService) GetMyCycles(ctx context.Context, userID string) ([]dto_response.GetSubscriptionCycleResponse, error) {
	return s.GetByUserID(ctx, userID)
}

func toCycleResponse(c entity.SubscriptionCycle) dto_response.GetSubscriptionCycleResponse {
	resp := dto_response.GetSubscriptionCycleResponse{
		ID: c.ID.String(), MembershipID: c.MembershipID.String(), CycleNumber: c.CycleNumber,
		PlanName: c.PlanName, PlanCategory: c.PlanCategory, DurationMonths: c.DurationMonths,
		AmountPaid: c.AmountPaid, PaymentChannel: c.PaymentChannel, TransactionID: c.TransactionID,
		Status: string(c.Status), StartDate: c.StartDate.Format("2006-01-02"), CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
	if c.PlanDurationID != nil { resp.PlanDurationID = c.PlanDurationID.String() }
	if c.EndDate != nil { e := c.EndDate.Format("2006-01-02"); resp.EndDate = &e }
	return resp
}
