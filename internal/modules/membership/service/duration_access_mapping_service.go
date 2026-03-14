package service

import (
	"context"
	"fmt"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/membership/repository"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	// "gorm.io/gorm"
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
		mappingRepository  repository.DurationAccessMappingRepository
		durationRepository repository.PlanDurationRepository
		entitlementRepo    repository.EntitlementRepository
	}
)

func NewDurationAccessMappingService(mappingRepo repository.DurationAccessMappingRepository, durationRepo repository.PlanDurationRepository, entitlementRepo repository.EntitlementRepository) DurationAccessMappingService {
	return &durationAccessMappingService{mappingRepository: mappingRepo, durationRepository: durationRepo, entitlementRepo: entitlementRepo}
}

func (s *durationAccessMappingService) Create(ctx context.Context, planID, planDurationID string, req dto_request.CreateDurationAccessMappingRequest) (dto_response.GetDurationAccessMappingResponse, error) {
	parsedPlanDurationID, _ := uuid.Parse(planDurationID)
	planDuration, err := s.durationRepository.GetByID(ctx, nil, parsedPlanDurationID)
	if err != nil || planDuration.PlanID.String() != planID { return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("plan duration") }

	parsedEntitlementID, _ := uuid.Parse(req.EntitlementID)
	entitlement, err := s.entitlementRepo.GetById(ctx, nil, parsedEntitlementID)
	if err != nil { return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("entitlement") }

	_, err = s.mappingRepository.GetByPlanDurationAndEntitlement(ctx, nil, parsedPlanDurationID, parsedEntitlementID)
	if err == nil { return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordAlreadyExist(fmt.Sprintf("mapping for '%s'", entitlement.Key)) }

	usageLimit := req.UsageLimit
	if entitlement.RestrictionType == entity.RestrictionFrequencyLimited && usageLimit <= 0 {
		return dto_response.GetDurationAccessMappingResponse{}, myerror.New("usage_limit > 0 required", myerror.Error_InvalidRequest)
	}

	newMapping := entity.NewDurationAccessMapping(parsedPlanDurationID, parsedEntitlementID, entitlement.Key, entitlement.Name, entitlement.Feature.Name, usageLimit)
	result, err := s.mappingRepository.Create(ctx, nil, newMapping)
	if err != nil { return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err) }
	result.Entitlement = entitlement
	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) GetAllByPlanDurationID(ctx context.Context, planID, planDurationID string) ([]dto_response.GetDurationAccessMappingResponse, error) {
	parsedID, _ := uuid.Parse(planDurationID)
	planDuration, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil || planDuration.PlanID.String() != planID { return nil, myerror.RecordNotFound("plan duration") }
	results, err := s.mappingRepository.GetAllByPlanDurationID(ctx, nil, parsedID)
	if err != nil { return nil, myerror.DatabaseError(err) }
	var res []dto_response.GetDurationAccessMappingResponse
	for _, r := range results { res = append(res, toMappingResponse(r)) }
	return res, nil
}

func (s *durationAccessMappingService) GetById(ctx context.Context, planID, planDurationID, id string) (dto_response.GetDurationAccessMappingResponse, error) {
	parsedID, _ := uuid.Parse(id)
	result, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil || result.PlanDurationID.String() != planDurationID || result.PlanDuration.PlanID.String() != planID { return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping") }
	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) Update(ctx context.Context, planID, planDurationID, id string, req dto_request.UpdateDurationAccessMappingRequest) (dto_response.GetDurationAccessMappingResponse, error) {
	parsedID, _ := uuid.Parse(id)
	existing, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil || existing.PlanDurationID.String() != planDurationID || existing.PlanDuration.PlanID.String() != planID { return dto_response.GetDurationAccessMappingResponse{}, myerror.RecordNotFound("access mapping") }
	existing.Status = entity.MappingStatus(req.Status)
	if req.UsageLimit != nil { existing.UsageLimit = *req.UsageLimit }
	result, err := s.mappingRepository.Update(ctx, nil, existing)
	if err != nil { return dto_response.GetDurationAccessMappingResponse{}, myerror.DatabaseError(err) }
	return toMappingResponse(result), nil
}

func (s *durationAccessMappingService) Delete(ctx context.Context, planID, planDurationID, id string) error {
	parsedID, _ := uuid.Parse(id)
	existing, err := s.mappingRepository.GetByID(ctx, nil, parsedID)
	if err != nil || existing.PlanDurationID.String() != planDurationID || existing.PlanDuration.PlanID.String() != planID { return myerror.RecordNotFound("access mapping") }
	return s.mappingRepository.Delete(ctx, nil, parsedID)
}

func toMappingResponse(m entity.DurationAccessMapping) dto_response.GetDurationAccessMappingResponse {
	return dto_response.GetDurationAccessMappingResponse{
		ID: m.ID.String(), PlanDurationID: m.PlanDurationID.String(), EntitlementID: m.EntitlementID.String(),
		EntitlementKey: m.EntitlementKey, EntitlementName: m.EntitlementName, Category: m.Category,
		RestrictionType: string(m.Entitlement.RestrictionType), TokenCost: m.Entitlement.TokenCost, UsageLimit: m.UsageLimit,
		Status: string(m.Status), CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
}

