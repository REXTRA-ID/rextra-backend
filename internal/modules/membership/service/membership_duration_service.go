package service

import (
	"context"
	"rextra-backend/internal/modules/membership/repository"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

type (
	MembershipDurationService interface {
		GetAllMembershipDuration(ctx context.Context) ([]dto_response.GetMembershipDurationResponse, error)
	}

	membershipDurationService struct {
		membershipDurationRepository repository.MembershipDurationRepository
		db                           *gorm.DB
	}
)

func NewMembershipDurationService(membershipDurationRepository repository.MembershipDurationRepository, db *gorm.DB) MembershipDurationService {
	return &membershipDurationService{
		membershipDurationRepository: membershipDurationRepository,
		db:                           db,
	}
}

func (s *membershipDurationService) GetAllMembershipDuration(ctx context.Context) ([]dto_response.GetMembershipDurationResponse, error) {
	membershipDuration, err := s.membershipDurationRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	return MapMembershipDurationResponse(membershipDuration), nil
}

func MapMembershipDurationResponse(membershipDuration []entity.MembershipDuration) []dto_response.GetMembershipDurationResponse {
	var res []dto_response.GetMembershipDurationResponse
	for i := range membershipDuration {
		data := dto_response.GetMembershipDurationResponse{
			ID:                   membershipDuration[i].ID.String(),
			DurationMonth:        membershipDuration[i].DurationMonth,
			TokenBonusPercentage: membershipDuration[i].TokenBonusPercentage,
			RextraPoinMultiplier: membershipDuration[i].RextraPoinMultiplier,
			IsActive:             membershipDuration[i].IsActive,
		}
		res = append(res, data)
	}
	return res
}
