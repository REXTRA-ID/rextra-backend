package service

import (
	"context"
	"rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
)

type (
	MembershipService interface {
		CheckExpiredMemberships() error
	}

	membershipService struct {
		membershipRepository     repository.MembershipRepository
		membershipPlanRepository repository.MembershipPlanRepository
	}
)

func NewMembershipService(membershipRepository repository.MembershipRepository,
	membershipPlanRepository repository.MembershipPlanRepository) MembershipService {
	return &membershipService{
		membershipRepository:     membershipRepository,
		membershipPlanRepository: membershipPlanRepository,
	}
}

func (s *membershipService) CheckExpiredMemberships() error {
	ctx := context.Background()

	planNonMember, err := s.membershipPlanRepository.GetByPlanName(ctx, nil, string(entity.PLANSTANDARD))
	if err != nil {
		return err
	}

	expiredMemberships, err := s.membershipRepository.GetExpiredMemberships(ctx, nil)
	if err != nil {
		return err
	}

	if len(expiredMemberships) == 0 {
		mylog.ColorizeInfo("no expired memberships")
		return nil
	}

	for _, membership := range expiredMemberships {
		membership.ExpiringMembership(&planNonMember)

		_, err = s.membershipRepository.Update(ctx, nil, membership)
		if err != nil {
			return err
		}
	}

	return nil
}
