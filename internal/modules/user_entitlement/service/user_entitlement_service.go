package service

import (
	"context"
	"errors"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/user_entitlement/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	UserEntitlementService interface {
		GetMyQuota(ctx context.Context, userID uuid.UUID, filter dto_request.MyQuotaFilter) (dto_response.MyQuotaResponse, error)
	}

	userEntitlementService struct {
		repo repository.UserEntitlementRepository
	}
)

func NewUserEntitlementService(repo repository.UserEntitlementRepository) UserEntitlementService {
	return &userEntitlementService{repo: repo}
}

func (s *userEntitlementService) GetMyQuota(ctx context.Context, userID uuid.UUID, filter dto_request.MyQuotaFilter) (dto_response.MyQuotaResponse, error) {
	membership, err := s.repo.GetActiveMembershipByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.MyQuotaResponse{PlanName: "", Entitlements: []dto_response.MyEntitlementQuotaResponse{}}, nil
		}
		return dto_response.MyQuotaResponse{}, myerror.DatabaseError(err)
	}

	if membership.DurationID == nil {
		return dto_response.MyQuotaResponse{PlanName: string(membership.PlanName), Entitlements: []dto_response.MyEntitlementQuotaResponse{}}, nil
	}

	mappings, err := s.repo.GetEntitlementsByPlanDurationID(ctx, membership.DurationID)
	if err != nil { return dto_response.MyQuotaResponse{}, myerror.DatabaseError(err) }

	quotas, err := s.repo.GetQuotasByMembershipID(ctx, membership.ID)
	if err != nil { return dto_response.MyQuotaResponse{}, myerror.DatabaseError(err) }

	quotaMap := make(map[string]entity.UserEntitlementQuota, len(quotas))
	for _, q := range quotas { quotaMap[q.EntitlementKey] = q }

	var entitlements []dto_response.MyEntitlementQuotaResponse
	for _, m := range mappings {
		e := m.Entitlement
		if filter.RestrictionType != nil && string(e.RestrictionType) != *filter.RestrictionType { continue }

		item := dto_response.MyEntitlementQuotaResponse{EntitlementKey: e.Key, EntitlementName: e.Name, RestrictionType: string(e.RestrictionType)}
		switch e.RestrictionType {
		case entity.RestrictionTokenGated:
			item.TokenCost = e.TokenCost
		case entity.RestrictionFrequencyLimited:
			if quota, ok := quotaMap[e.Key]; ok {
				granted, used, remaining := quota.QuotaGranted, quota.QuotaUsed, quota.QuotaRemaining
				item.QuotaGranted, item.QuotaUsed, item.QuotaRemaining, item.CycleExpiredAt = &granted, &used, &remaining, quota.CycleExpiredAt
			}
		}
		entitlements = append(entitlements, item)
	}

	if entitlements == nil { entitlements = []dto_response.MyEntitlementQuotaResponse{} }
	return dto_response.MyQuotaResponse{PlanName: string(membership.PlanName), Entitlements: entitlements}, nil
}
