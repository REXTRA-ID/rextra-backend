package service

import (
	"context"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/hak_akses/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	HakAksesService interface {
		MakeHakAkses(ctx context.Context, dto dto_request.MakeHakAksesRequest) (dto_response.MakeHakAksesResponse, error)
		CheckAccess(ctx context.Context, dto dto_request.CheckAccessRequest) (bool, error)
	}

	hakAksesService struct {
		hakAksesRepository repository.HakAksesRepository
		db                 *gorm.DB
	}
)

func NewHakAksesService(hakAksesRepo repository.HakAksesRepository, db *gorm.DB) HakAksesService {
	return &hakAksesService{
		hakAksesRepository: hakAksesRepo,
		db:                 db,
	}
}

func (s *hakAksesService) MakeHakAkses(ctx context.Context, dto dto_request.MakeHakAksesRequest) (dto_response.MakeHakAksesResponse, error) {

	hakAkses := entity.NewHakAkses(uuid.MustParse(dto.MembershipPlanId), dto.Feature, dto.Action)

	newHakAkses, err := s.hakAksesRepository.Create(ctx, nil, hakAkses)
	if err != nil {
		return dto_response.MakeHakAksesResponse{}, err
	}

	return dto_response.MakeHakAksesResponse{
		Id:               newHakAkses.ID.String(),
		MembershipPlanId: newHakAkses.MembershipPlanId.String(),
		Feature:          string(newHakAkses.Feature),
		Action:           string(newHakAkses.Action),
	}, nil
}

func (s *hakAksesService) CheckAccess(ctx context.Context, dto dto_request.CheckAccessRequest) (bool, error) {

	hakAkses, err := s.hakAksesRepository.FindByFeatureAction(ctx, nil, dto.Feature, dto.Action)
	if err != nil {
		return false, err
	}

	if dto.MembershipPlan == string(hakAkses.MembershipPlan.PlanName) {
		return true, nil
	}

	return false, nil
}
