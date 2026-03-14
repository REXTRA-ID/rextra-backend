package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		GetStarterPlan(ctx context.Context) (entity.MembershipPlans, error)
	}

	membershipPlanService struct {
		planRepository repository.MembershipPlanRepository
	}
)

func NewMembershipPlanService(planRepo repository.MembershipPlanRepository) MembershipPlanService {
	return &membershipPlanService{planRepository: planRepo}
}

func (s *membershipPlanService) Create(ctx context.Context, req dto_request.CreateMembershipPlanRequest) (dto_response.GetMembershipPlanResponse, error) {
	_, err := s.planRepository.GetByPlanName(ctx, nil, entity.PlanName(req.PlanName))
	if err == nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.RecordAlreadyExist(fmt.Sprintf("membership plan '%s'", req.PlanName))
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}

	benefitsJSON, _ := json.Marshal(req.Benefits)
	newPlan := entity.NewMembershipPlan(
		entity.PlanName(req.PlanName),
		entity.PlanCategory(req.Category),
		req.TierLabel,
		req.EmblemKey,
		req.Description,
		entity.PricingMode(req.PricingMode),
		entity.DurationMode(req.DurationMode),
		req.BasePrice1M,
		req.BaseToken1M,
	)
	newPlan.Discount3M = req.Discount3M
	newPlan.Discount6M = req.Discount6M
	newPlan.Discount12M = req.Discount12M
	newPlan.BonusToken3M = req.BonusToken3M
	newPlan.BonusToken6M = req.BonusToken6M
	newPlan.BonusToken12M = req.BonusToken12M
	newPlan.Benefits = datatypes.JSON(benefitsJSON)

	result, err := s.planRepository.Create(ctx, nil, newPlan)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}
	return toPlanResponse(result, nil), nil
}

func (s *membershipPlanService) GetAll(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error) {
	results, err := s.planRepository.GetAllWithDurations(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}
	var res []dto_response.GetMembershipPlanResponse
	for _, r := range results {
		res = append(res, toPlanResponse(r, r.PlanDurations))
	}
	return res, nil
}

func (s *membershipPlanService) GetById(ctx context.Context, id string) (dto_response.GetMembershipPlanResponse, error) {
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
	return toPlanResponse(result, result.PlanDurations), nil
}

func (s *membershipPlanService) Update(ctx context.Context, id string, req dto_request.UpdateMembershipPlanRequest) (dto_response.GetMembershipPlanResponse, error) {
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

	benefitsJSON, _ := json.Marshal(req.Benefits)
	existing.TierLabel = req.TierLabel
	existing.EmblemKey = req.EmblemKey
	existing.Description = req.Description
	existing.Status = entity.PlanStatus(req.Status)
	existing.PricingMode = entity.PricingMode(req.PricingMode)
	existing.BasePrice1M = req.BasePrice1M
	existing.BaseToken1M = req.BaseToken1M
	existing.Discount3M = req.Discount3M
	existing.Discount6M = req.Discount6M
	existing.Discount12M = req.Discount12M
	existing.BonusToken3M = req.BonusToken3M
	existing.BonusToken6M = req.BonusToken6M
	existing.BonusToken12M = req.BonusToken12M
	existing.Benefits = datatypes.JSON(benefitsJSON)
	existing.UpdatedAt = time.Now().UTC()

	if existing.PlanName == entity.PlanNameStarter && req.StarterDurationMonths > 0 {
		existing.StarterDurationMonths = req.StarterDurationMonths
	}

	result, err := s.planRepository.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.GetMembershipPlanResponse{}, myerror.DatabaseError(err)
	}
	return toPlanResponse(result, nil), nil
}

func (s *membershipPlanService) Delete(ctx context.Context, id string) error {
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
	count, err := s.planRepository.CountActiveMembers(ctx, nil, parsedID)
	if err != nil {
		return myerror.DatabaseError(err)
	}
	if count > 0 {
		return myerror.New(fmt.Sprintf("plan cannot be deleted because it has %d active member(s)", count), myerror.Error_InvalidRequest)
	}
	return s.planRepository.Delete(ctx, nil, parsedID)
}

func (s *membershipPlanService) GetCatalog(ctx context.Context) ([]dto_response.GetMembershipPlanResponse, error) {
	results, err := s.planRepository.GetAllWithDurations(ctx, nil)
	if err != nil {
		return nil, myerror.DatabaseError(err)
	}
	var res []dto_response.GetMembershipPlanResponse
	for _, r := range results {
		if r.Status == entity.PlanStatusActive {
			resp := toPlanResponse(r, r.PlanDurations)
			resp.ActiveUsers = 0
			res = append(res, resp)
		}
	}
	return res, nil
}

func (s *membershipPlanService) GetStarterPlan(ctx context.Context) (entity.MembershipPlans, error) {
	plan, err := s.planRepository.GetByPlanName(ctx, nil, entity.PlanNameStarter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.MembershipPlans{}, myerror.RecordNotFound("starter plan")
		}
		return entity.MembershipPlans{}, myerror.DatabaseError(err)
	}
	return plan, nil
}

func toPlanResponse(p entity.MembershipPlans, durations []entity.PlanDuration) dto_response.GetMembershipPlanResponse {
	benefits, _ := p.GetBenefits()
	if benefits == nil { benefits = []string{} }
	createdAt, updatedAt := "", ""
	if !p.CreatedAt.IsZero() { createdAt = p.CreatedAt.Format(time.RFC3339) }
	if !p.UpdatedAt.IsZero() { updatedAt = p.UpdatedAt.Format(time.RFC3339) }

	res := dto_response.GetMembershipPlanResponse{
		ID: p.ID.String(),
		PlanName: string(p.PlanName),
		Category: string(p.Category),
		TierLabel: p.TierLabel,
		EmblemKey: p.EmblemKey,
		Description: p.Description,
		Status: string(p.Status),
		PricingMode: string(p.PricingMode),
		DurationMode: string(p.DurationMode),
		BasePrice1M: p.BasePrice1M,
		BaseToken1M: p.BaseToken1M,
		Discount3M: p.Discount3M,
		Discount6M: p.Discount6M,
		Discount12M: p.Discount12M,
		BonusToken3M: p.BonusToken3M,
		BonusToken6M: p.BonusToken6M,
		BonusToken12M: p.BonusToken12M,
		ActiveUsers: p.ActiveUsers,
		Benefits: benefits,
		StarterDurationMonths: p.StarterDurationMonths,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	for _, d := range durations {
		res.PlanDurations = append(res.PlanDurations, toDurationResponse(d, 0))
	}
	return res
}
