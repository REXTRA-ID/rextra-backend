package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/membership/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	MembershipPlanService interface {
		Create(ctx context.Context, req dto_request.CreateMembershipPlanRequest) (dto_response.GetMembershipPlanResponse, error)
		GetAll(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetMembershipPlanResponse, error)
		Update(ctx context.Context, id string, req dto_request.UpdateMembershipPlanRequest) (dto_response.GetMembershipPlanResponse, error)
		Delete(ctx context.Context, id string) error

		GetCatalog(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error)
	}

	membershipPlanService struct {
		planRepository repository.MembershipPlanRepository
	}
)

func NewMembershipPlanService(planRepo repository.MembershipPlanRepository) MembershipPlanService {
	return &membershipPlanService{
		planRepository: planRepo,
	}
}

func (s *membershipPlanService) Create(
	ctx context.Context,
	req dto_request.CreateMembershipPlanRequest,
) (dto_response.GetMembershipPlanResponse, error) {

	// cek duplicate plan name
	_, err := s.planRepository.GetByPlanName(ctx, nil, entity.EnumPlanName(req.PlanName))
	if err == nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.RecordAlreadyExist("membership plan")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	benefitsJSON, err := json.Marshal(req.Benefits)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.ProcessingError(err)
	}

	newPlan := entity.MembershipPlans{
		PlanName:         entity.EnumPlanName(req.PlanName),
		MonthlyToken:     req.MonthlyToken,
		BaseMonthlyPrice: req.BaseMonthlyPrice,
		Description:      req.Description,
		Benefits:         datatypes.JSON(benefitsJSON),
		IsActive:         true,
	}

	result, err := s.planRepository.Create(ctx, nil, newPlan)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	return toPlanResponse(result), nil
}

func (s *membershipPlanService) GetAll(
	ctx context.Context,
) ([]dto_response.GetMembershipPlanResponse, error) {

	results, err := s.planRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetMembershipPlanResponse

	for _, r := range results {
		res = append(res, toPlanResponse(r))
	}

	return res, nil
}

func (s *membershipPlanService) GetById(
	ctx context.Context,
	id string,
) (dto_response.GetMembershipPlanResponse, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.planRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetMembershipPlanResponse{}, myerror.RecordNotFound("membership plan")
		}
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	return toPlanResponse(result), nil
}

func (s *membershipPlanService) Update(
	ctx context.Context,
	id string,
	req dto_request.UpdateMembershipPlanRequest,
) (dto_response.GetMembershipPlanResponse, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.planRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetMembershipPlanResponse{}, myerror.RecordNotFound("membership plan")
		}
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	benefitsJSON, err := json.Marshal(req.Benefits)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.ProcessingError(err)
	}

	existing.PlanName = entity.EnumPlanName(req.PlanName)
	existing.MonthlyToken = req.MonthlyToken
	existing.BaseMonthlyPrice = req.BaseMonthlyPrice
	existing.Description = req.Description
	existing.Benefits = datatypes.JSON(benefitsJSON)
	existing.IsActive = req.IsActive
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.planRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	return toPlanResponse(result), nil
}

func (s *membershipPlanService) Delete(
	ctx context.Context,
	id string,
) error {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	_, err = s.planRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("membership plan")
		}
		return myerror.DatabaseError(err)
	}

	return s.planRepository.Delete(ctx, nil, parsedID)
}

func toPlanResponse(p entity.MembershipPlans) dto_response.GetMembershipPlanResponse {

	benefits, _ := p.GetBenefits()
	if benefits == nil {
		benefits = []string{}
	}

	createdAt := ""
	updatedAt := ""

	if !p.CreatedAt.IsZero() {
		createdAt = p.CreatedAt.Format(time.RFC3339)
	}

	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format(time.RFC3339)
	}

	return dto_response.GetMembershipPlanResponse{
		ID:               p.ID.String(),
		PlanName:         string(p.PlanName),
		MonthlyToken:     p.MonthlyToken,
		BaseMonthlyPrice: p.BaseMonthlyPrice,
		Description:      p.Description,
		Benefits:         benefits,
		IsActive:         p.IsActive,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (s *membershipPlanService) GetCatalog(
	ctx context.Context,
) ([]dto_response.GetMembershipPlanResponse, error) {

	results, err := s.planRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetMembershipPlanResponse

	for _, r := range results {

		if !r.IsActive {
			continue
		}

		res = append(res, toPlanResponse(r))
	}

	return res, nil
}
