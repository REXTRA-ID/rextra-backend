package service

import (
	"context"
	"errors"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/membership/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	MembershipDurationService interface {
		Create(ctx context.Context, req dto_request.CreateMembershipDurationRequest) (dto_response.GetMembershipDurationResponse, error)
		GetAll(ctx context.Context, planID string) ([]dto_response.GetMembershipDurationResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetMembershipDurationResponse, error)
		Update(ctx context.Context, id string, req dto_request.UpdateMembershipDurationRequest) (dto_response.GetMembershipDurationResponse, error)
		Delete(ctx context.Context, id string) error
	}

	membershipDurationService struct {
		durationRepository repository.MembershipDurationRepository
	}
)

func NewMembershipDurationService(
	durationRepo repository.MembershipDurationRepository,
) MembershipDurationService {
	return &membershipDurationService{
		durationRepository: durationRepo,
	}
}

func (s *membershipDurationService) Create(
	ctx context.Context,
	req dto_request.CreateMembershipDurationRequest,
) (dto_response.GetMembershipDurationResponse, error) {

	newDuration := entity.NewMembershipDuration(
		req.DurationMonth,
		req.TokenBonusPercentage,
		int(req.RextraPoinMultiplier),
	)

	result, err := s.durationRepository.Create(ctx, nil, newDuration)
	if err != nil {
		return dto_response.GetMembershipDurationResponse{}, myerror.DatabaseError(err)
	}

	return toMembershipDurationResponse(result), nil
}

func (s *membershipDurationService) GetAll(
	ctx context.Context,
	planID string,
) ([]dto_response.GetMembershipDurationResponse, error) {

	parsedPlanID, err := uuid.Parse(planID)
	if err != nil {
		return nil, myerror.InvalidRequest(err)
	}

	results, err := s.durationRepository.GetAllByPlanID(ctx, nil, parsedPlanID)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetMembershipDurationResponse

	for _, r := range results {
		res = append(res, toMembershipDurationResponse(r))
	}

	return res, nil
}

func (s *membershipDurationService) GetById(
	ctx context.Context,
	id string,
) (dto_response.GetMembershipDurationResponse, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetMembershipDurationResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetMembershipDurationResponse{}, myerror.RecordNotFound("membership duration")
		}
		return dto_response.GetMembershipDurationResponse{}, myerror.DatabaseError(err)
	}

	return toMembershipDurationResponse(result), nil
}

func (s *membershipDurationService) Update(
	ctx context.Context,
	id string,
	req dto_request.UpdateMembershipDurationRequest,
) (dto_response.GetMembershipDurationResponse, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetMembershipDurationResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetMembershipDurationResponse{}, myerror.RecordNotFound("membership duration")
		}
		return dto_response.GetMembershipDurationResponse{}, myerror.DatabaseError(err)
	}

	// DurationMonth tidak boleh diubah
	existing.TokenBonusPercentage = req.TokenBonusPercentage
	existing.RextraPoinMultiplier = int(req.RextraPoinMultiplier)
	existing.IsActive = req.IsActive
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.durationRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetMembershipDurationResponse{}, myerror.DatabaseError(err)
	}

	return toMembershipDurationResponse(result), nil
}

func (s *membershipDurationService) Delete(
	ctx context.Context,
	id string,
) error {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	_, err = s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("membership duration")
		}
		return myerror.DatabaseError(err)
	}

	return s.durationRepository.Delete(ctx, nil, parsedID)
}

func toMembershipDurationResponse(
	d entity.MembershipDuration,
) dto_response.GetMembershipDurationResponse {

	createdAt := ""
	updatedAt := ""

	if !d.CreatedAt.IsZero() {
		createdAt = d.CreatedAt.Format(time.RFC3339)
	}

	if !d.UpdatedAt.IsZero() {
		updatedAt = d.UpdatedAt.Format(time.RFC3339)
	}

	return dto_response.GetMembershipDurationResponse{
		ID:                   d.ID.String(),
		DurationMonth:        d.DurationMonth,
		TokenBonusPercentage: d.TokenBonusPercentage,
		RextraPoinMultiplier: d.RextraPoinMultiplier,
		IsActive:             d.IsActive,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}
