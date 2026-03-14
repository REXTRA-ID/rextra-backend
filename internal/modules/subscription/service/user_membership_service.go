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
	UserMembershipService interface {
		GetAll(ctx context.Context, req dto_request.UserMembershipFilterRequest) (dto_response.PaginatedUserMembershipResponse, error)
		GetById(ctx context.Context, membershipID string) (dto_response.GetUserMembershipResponse, error)
		GetMyMembership(ctx context.Context, userID string) (dto_response.GetUserMembershipResponse, error)
	}

	userMembershipService struct {
		membershipRepo membershipRepo.UserMembershipRepository
		userRepo       repository.UserRepository
	}
)

func NewUserMembershipService(membershipRepo membershipRepo.UserMembershipRepository, userRepo repository.UserRepository) UserMembershipService {
	return &userMembershipService{membershipRepo: membershipRepo, userRepo: userRepo}
}

func (s *userMembershipService) GetAll(ctx context.Context, req dto_request.UserMembershipFilterRequest) (dto_response.PaginatedUserMembershipResponse, error) {
	page, pageSize := req.Page, req.PageSize
	if pageSize <= 0 { pageSize = 20 }
	offset := page * pageSize
	memberships, total, err := s.membershipRepo.GetAllPaginated(ctx, nil, offset, pageSize, req.PlanName, req.Search, req.IsActive)
	if err != nil { return dto_response.PaginatedUserMembershipResponse{}, myerror.DatabaseError(err) }
	userIDs := make([]uuid.UUID, len(memberships))
	for i, m := range memberships { userIDs[i] = m.UserID }
	userMap, err := s.userRepo.GetByIDs(ctx, nil, userIDs)
	if err != nil { return dto_response.PaginatedUserMembershipResponse{}, myerror.DatabaseError(err) }
	var data []dto_response.GetUserMembershipListResponse
	for _, m := range memberships {
		user := userMap[m.UserID]
		data = append(data, toListResponse(m, user.Fullname, user.Email))
	}
	return dto_response.PaginatedUserMembershipResponse{Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize)))}, nil
}

func (s *userMembershipService) GetById(ctx context.Context, membershipID string) (dto_response.GetUserMembershipResponse, error) {
	parsedID, err := uuid.Parse(membershipID)
	if err != nil { return dto_response.GetUserMembershipResponse{}, myerror.InvalidRequest(err) }
	membership, err := s.membershipRepo.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetUserMembershipResponse{}, myerror.RecordNotFound("membership") }
		return dto_response.GetUserMembershipResponse{}, myerror.DatabaseError(err)
	}
	return toDetailResponse(membership), nil
}

func (s *userMembershipService) GetMyMembership(ctx context.Context, userID string) (dto_response.GetUserMembershipResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return dto_response.GetUserMembershipResponse{}, myerror.InvalidRequest(err) }
	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetUserMembershipResponse{}, myerror.RecordNotFound("membership") }
		return dto_response.GetUserMembershipResponse{}, myerror.DatabaseError(err)
	}
	return toDetailResponse(membership), nil
}

func toDetailResponse(m entity.Memberships) dto_response.GetUserMembershipResponse {
	resp := dto_response.GetUserMembershipResponse{
		MembershipID: m.ID.String(), UserID: m.UserID.String(), PlanName: string(m.PlanName),
		IsActive: m.IsActive, AutoRenew: m.AutoRenew, CurrentTokenBalance: m.CurrentTokenBalance,
		CurrentPoinBalance: m.CurrentPoinBalance, PaidCycleCount: m.PaidCycleCount,
		EntitlementCount: m.EntitlementCount, RemainingDays: m.RemainingDays(), UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
	if m.PlanID != nil { resp.PlanID = m.PlanID.String() }
	if m.DurationID != nil { resp.DurationID = m.DurationID.String() }
	if m.DurationMonths != nil { resp.DurationMonths = m.DurationMonths }
	if m.StartedAt != nil { s := m.StartedAt.Format(time.RFC3339); resp.StartedAt = &s }
	if m.ExpiredAt != nil { e := m.ExpiredAt.Format(time.RFC3339); resp.ExpiredAt = &e }
	return resp
}

func toListResponse(m entity.Memberships, userName, userEmail string) dto_response.GetUserMembershipListResponse {
	resp := dto_response.GetUserMembershipListResponse{
		MembershipID: m.ID.String(), UserID: m.UserID.String(), UserName: userName, UserEmail: userEmail,
		PlanName: string(m.PlanName), DurationMonths: m.DurationMonths, IsActive: m.IsActive,
		RemainingDays: m.RemainingDays(), PaidCycleCount: m.PaidCycleCount, UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
	if m.ExpiredAt != nil { e := m.ExpiredAt.Format(time.RFC3339); resp.ExpiredAt = &e }
	return resp
}
