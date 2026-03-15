package service

import (
	"context"
	"rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/entity"
)

type (
	MembershipService interface {
		CheckExpiredMemberships() error
	}

	membershipService struct {
		membershipRepository     repository.UserMembershipRepository
		membershipPlanRepository repository.MembershipPlanRepository
	}
)

func NewMembershipService(membershipRepository repository.UserMembershipRepository,
	membershipPlanRepository repository.MembershipPlanRepository) MembershipService {
	return &membershipService{
		membershipRepository:     membershipRepository,
		membershipPlanRepository: membershipPlanRepository,
	}
}

func (s *membershipService) CheckExpiredMemberships() error {
	ctx := context.Background()

	_, err := s.membershipPlanRepository.GetByPlanName(ctx, nil, entity.PlanName(entity.PLANSTANDARD))
	if err != nil {
		return err
	}

	return nil
}
