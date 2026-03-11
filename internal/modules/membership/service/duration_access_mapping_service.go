package service

import (
	"context"
	"errors"
	"fmt"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	e "rextra-backend/internal/modules/entitlement/repository"
	m "rextra-backend/internal/modules/membership/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	DurationAccessMappingService interface {
		Create(ctx context.Context, planID, planDurationID string, req dto_request.CreateDurationAccessMappingRequest) (dto_response.GetDurationAccessMappingResponse, error)
		GetAllByPlanDurationID(ctx context.Context, planID, planDurationID string) ([]dto_response.GetDurationAccessMappingResponse, error)
		GetById(ctx context.Context, planID, planDurationID, id string) (dto_response.GetDurationAccessMappingResponse, error)
		Update(ctx context.Context, planID, planDurationID, id string, req dto_request.UpdateDurationAccessMappingRequest) (dto_response.GetDurationAccessMappingResponse, error)
		Delete(ctx context.Context, planID, planDurationID, id string) error
	}

	durationAccessMappingService struct {
		mappingRepository  m.DurationAccessMappingRepository
		durationRepository m.MembershipDurationRepository
		entitlementRepo    e.EntitlementRepository
	}
)

func NewDurationAccessMappingService(
	mappingRepo m.DurationAccessMappingRepository,
	durationRepo m.MembershipDurationRepository,
	entitlementRepo e.EntitlementRepository,
) DurationAccessMappingService {
	return &durationAccessMappingService{
		mappingRepository:  mappingRepo,
		durationRepository: durationRepo,
		entitlementRepo:    entitlementRepo,
	}
}

func (s *durationAccessMappingService) Create(
	ctx context.Context,
	planID, planDurationID string,
	req dto_request.CreateDurationAccessMappingRequest,
) (dto_response.GetDurationAccessMappingResponse, error) {

	parsedPlanDurationID, err := uuid.Parse(planDurationID)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.InvalidRequest(err)
	}

	// pastikan duration ada
	_, err = s.durationRepository.GetByID(ctx, nil, parsedPlanDurationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("plan duration")
		}
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	parsedEntitlementID, err := uuid.Parse(req.EntitlementID)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.InvalidRequest(err)
	}

	entitlement, err := s.entitlementRepo.GetById(ctx, nil, parsedEntitlementID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("entitlement")
		}
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	// cek duplicate
	_, err = s.mappingRepository.GetByPlanDurationAndEntitlement(ctx, nil, parsedPlanDurationID, parsedEntitlementID)
	if err == nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordAlreadyExist(
			fmt.Sprintf("mapping for entitlement '%s' already exists", entitlement.Key),
		)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	category := entitlement.Feature.Name

	newMapping := entity.NewDurationAccessMapping(
		parsedPlanDurationID,
		parsedEntitlementID,
		entitlement.Key,
		entitlement.Name,
		category,
		req.UsageLimit,
	)

	result, err := s.mappingRepository.Create(ctx, nil, newMapping)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) GetAllByPlanDurationID(
	ctx context.Context,
	planID, planDurationID string,
) ([]dto_response.GetDurationAccessMappingResponse, error) {

	parsedID, err := uuid.Parse(planDurationID)
	if err != nil {
		return nil, myerror.InvalidRequest(err)
	}

	// hanya cek duration ada
	_, err = s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerror.RecordNotFound("plan duration")
		}
		return nil, myerror.DatabaseError(err)
	}

	results, err := s.mappingRepository.GetAllByPlanDurationID(ctx, nil, parsedID)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}

	var res []dto_response.GetDurationAccessMappingResponse
	for _, r := range results {
		res = append(res, toMappingResponse(r))
	}

	return res, nil
}

func (s *durationAccessMappingService) GetById(ctx context.Context, planID, planDurationID, id string) (dto_response.GetDurationAccessMappingResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.InvalidRequest(err)
	}

	result, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping")
		}
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	// Ownership check — dua level: mapping milik durationId yang benar, dan duration milik planId yang benar
	if result.PlanDurationID.String() != planDurationID {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping")
	}

	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) Update(ctx context.Context, planID, planDurationID, id string, req dto_request.UpdateDurationAccessMappingRequest) (dto_response.GetDurationAccessMappingResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.InvalidRequest(err)
	}

	existing, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping")
		}
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	// Ownership check — dua level
	if existing.PlanDurationID.String() != planDurationID {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping")
	}

	// Validasi restriction-specific fields
	if entity.RestrictionType(req.RestrictionType) == entity.RestrictionTokenGated && req.TokenCost <= 0 {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.New(
			"token_cost must be greater than 0 when restriction_type is 'token_gated'",
			myerror.Error_InvalidRequest,
		)
	}
	if entity.RestrictionType(req.RestrictionType) == entity.RestrictionFrequencyLimited {
		if req.UsageLimit <= 0 {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.New(
				"usage_limit must be greater than 0 when restriction_type is 'frequency_limited'",
				myerror.Error_InvalidRequest,
			)
		}
		if req.ResetPeriod == nil || *req.ResetPeriod == "" {
			return dto_response.GetDurationAccessMappingResponse{}, myerror.New(
				"reset_period is required when restriction_type is 'frequency_limited'",
				myerror.Error_InvalidRequest,
			)
		}
	}

	existing.UsageLimit = req.UsageLimit
	existing.Status = entity.MappingStatus(req.Status)

	result, err := s.mappingRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err)
	}

	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) Delete(
	ctx context.Context,
	planID, planDurationID, id string,
) error {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return myerror.InvalidRequest(err)
	}

	existing, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.RecordNotFound("access mapping")
		}
		return myerror.DatabaseError(err)
	}

	// ownership check
	if existing.PlanDurationID.String() != planDurationID {
		return myerror.RecordNotFound("access mapping")
	}

	return s.mappingRepository.Delete(ctx, nil, parsedID)
}

func toMappingResponse(m entity.DurationAccessMapping) dto_response.GetDurationAccessMappingResponse {
	return dto_response.GetDurationAccessMappingResponse{
		ID:              m.ID.String(),
		PlanDurationID:  m.PlanDurationID.String(),
		EntitlementID:   m.EntitlementID.String(),
		EntitlementKey:  m.EntitlementKey,
		EntitlementName: m.EntitlementName,
		Category:        m.Category,
		UsageLimit:      m.UsageLimit,
		Status:          string(m.Status),
		CreatedAt:       m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
