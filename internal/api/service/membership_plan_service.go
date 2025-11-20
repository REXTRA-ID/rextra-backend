package service

import (
	"context"
	"rextra-backend/internal/api/repository"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	MembershipPlanService interface {
		GetAllMembershipPlan(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error)
	}

	membershipPlanService struct {
		membershipPlanRepository repository.MembershipPlanRepository
		db                       *gorm.DB
	}
)

func NewMembershipPlanService(membershipRepository repository.MembershipPlanRepository, db *gorm.DB) MembershipPlanService {
	return &membershipPlanService{
		membershipPlanRepository: membershipRepository,
		db:                       db,
	}
}

func (s *membershipPlanService) GetAllMembershipPlan(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error) {
	plans, err := s.membershipPlanRepository.GetAll(ctx, nil)
	if err != nil {
		return []dto_response.GetMembershipPlanResponse{}, err
	}

	return MapMembershipPlanResponse(plans), nil
}

func MapMembershipPlanResponse(plans []entity.MembershipPlans) []dto_response.GetMembershipPlanResponse {
	var res []dto_response.GetMembershipPlanResponse

	for _, plan := range plans {
		res = append(res,
			dto_response.GetMembershipPlanResponse{
				ID:               plan.ID.String(),
				PlanName:         string(plan.PlanName),
				MonthlyToken:     plan.MonthlyToken,
				BaseMonthlyPrice: plan.BaseMonthlyPrice,
				Description:      plan.Description,
				Benefits:         plan.Benefits,
				IsActive:         plan.IsActive,
			},
		)
	}
	return res
}
