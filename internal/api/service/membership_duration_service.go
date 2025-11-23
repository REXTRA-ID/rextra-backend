package service

import (
	"context"
	"rextra-backend/internal/api/repository"
)

type (
	MembershipDurationService interface {
		GetAllMembershipDuration(ctx context.Context)
	}

	membershipDurationService struct {
		membershipDurationRepository repository.MembershipDurationRepository
	}
)

func NewMembershipDurationService(membershipDurationRepository repository.MembershipDurationRepository) MembershipDurationService {
	return &membershipDurationService{
		membershipDurationRepository: membershipDurationRepository,
	}
}

func (m *membershipDurationService) GetAllMembershipDuration(ctx context.Context) {

}
