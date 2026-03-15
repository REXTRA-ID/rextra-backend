package service

import (
	"context"
	"errors"
	"fmt"
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
	PlanDurationService interface {
		Create(ctx context.Context, planID string, req dto_request.CreatePlanDurationRequest) (dto_response.GetPlanDurationResponse, error)
		GetAllByPlanID(ctx context.Context, planID string) ([]dto_response.GetPlanDurationResponse, error)
		GetById(ctx context.Context, planID, id string) (dto_response.GetPlanDurationResponse, error)
		Update(ctx context.Context, planID, id string, req dto_request.UpdatePlanDurationRequest) (dto_response.GetPlanDurationResponse, error)
		Delete(ctx context.Context, planID, id string) error
	}

	planDurationService struct {
		planRepository     repository.MembershipPlanRepository
		durationRepository repository.PlanDurationRepository
	}
)

func NewPlanDurationService(planRepo repository.MembershipPlanRepository, durationRepo repository.PlanDurationRepository) PlanDurationService {
	return &planDurationService{planRepository: planRepo, durationRepository: durationRepo}
}

func (s *planDurationService) Create(ctx context.Context, planID string, req dto_request.CreatePlanDurationRequest) (dto_response.GetPlanDurationResponse, error) {
	parsedPlanID, err := uuid.Parse(planID)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.InvalidRequest(err) }

	plan, err := s.planRepository.GetByID(ctx, nil, parsedPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetPlanDurationResponse{}, myerror.RecordNotFound("membership plan") }
		return dto_response.GetPlanDurationResponse{}, myerror.DatabaseError(err)
	}

	if plan.DurationMode == entity.DurationModeWithoutDuration {
		return dto_response.GetPlanDurationResponse{}, myerror.New(fmt.Sprintf("plan '%s' is 'tanpa_durasi'", plan.PlanName), myerror.Error_InvalidRequest)
	}

	if plan.PricingMode == entity.PricingModeManual && req.FinalPrice <= 0 {
		return dto_response.GetPlanDurationResponse{}, myerror.New("final_price is required for manual pricing", myerror.Error_InvalidRequest)
	}

	_, err = s.durationRepository.GetByPlanAndMonths(ctx, nil, parsedPlanID, req.DurationMonths)
	if err == nil { return dto_response.GetPlanDurationResponse{}, myerror.RecordAlreadyExist(fmt.Sprintf("duration %d months", req.DurationMonths)) }

	finalPrice := req.FinalPrice
	discountPct := req.DiscountPct
	if plan.PricingMode == entity.PricingModeAutomatic {
		discountMap := map[int]float64{1: 0, 3: plan.Discount3M, 6: plan.Discount6M, 12: plan.Discount12M}
		discountPct = discountMap[req.DurationMonths]
		finalPrice = int64(float64(plan.BasePrice1M) * float64(req.DurationMonths) * (1 - discountPct/100))
		bonusTokenMap := map[int]int{1: 0, 3: plan.BonusToken3M, 6: plan.BonusToken6M, 12: plan.BonusToken12M}
		req.TokenAmount = plan.BaseToken1M*req.DurationMonths + bonusTokenMap[req.DurationMonths]
		req.BonusToken = bonusTokenMap[req.DurationMonths]
	}

	if plan.Category == entity.PlanCategoryPaid && req.DurationPrice >= finalPrice {
		return dto_response.GetPlanDurationResponse{}, myerror.New("duration_price must be less than final_price", myerror.Error_InvalidRequest)
	}

	newDuration := entity.NewPlanDuration(parsedPlanID, req.DurationMonths, req.Price, finalPrice, req.DurationPrice, discountPct, req.TokenAmount, req.BonusToken)
	newDuration.PointsActive = req.PointsActive
	newDuration.PointsValue = req.PointsValue
	newDuration.BonusPoints = req.BonusPoints

	result, err := s.durationRepository.Create(ctx, nil, newDuration)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.DatabaseError(err) }
	return toDurationResponse(result, 0), nil
}

func (s *planDurationService) GetAllByPlanID(ctx context.Context, planID string) ([]dto_response.GetPlanDurationResponse, error) {
	parsedPlanID, err := uuid.Parse(planID)
	if err != nil { return nil, myerror.InvalidRequest(err) }
	_, err = s.planRepository.GetByID(ctx, nil, parsedPlanID)
	if err != nil { return nil, myerror.RecordNotFound("membership plan") }

	results, countMap, err := s.durationRepository.GetAllByPlanIDWithMappingCounts(ctx, nil, parsedPlanID)
	if err != nil { return nil, myerror.DatabaseError(err) }
	var res []dto_response.GetPlanDurationResponse
	for _, r := range results { res = append(res, toDurationResponse(r, int(countMap[r.ID]))) }
	return res, nil
}

func (s *planDurationService) GetById(ctx context.Context, planID, id string) (dto_response.GetPlanDurationResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.InvalidRequest(err) }
	result, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.RecordNotFound("plan duration") }
	if result.PlanID.String() != planID { return dto_response.GetPlanDurationResponse{}, myerror.RecordNotFound("plan duration") }
	count, _ := s.durationRepository.CountMappings(ctx, nil, result.ID)
	return toDurationResponse(result, int(count)), nil
}

func (s *planDurationService) Update(ctx context.Context, planID, id string, req dto_request.UpdatePlanDurationRequest) (dto_response.GetPlanDurationResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.InvalidRequest(err) }
	existing, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.RecordNotFound("plan duration") }
	if existing.PlanID.String() != planID { return dto_response.GetPlanDurationResponse{}, myerror.RecordNotFound("plan duration") }

	if existing.Plan.Category == entity.PlanCategoryPaid && req.DurationPrice >= req.FinalPrice {
		return dto_response.GetPlanDurationResponse{}, myerror.New("duration_price must be less than final_price", myerror.Error_InvalidRequest)
	}

	existing.Price = req.Price
	existing.DiscountPct = req.DiscountPct
	existing.FinalPrice = req.FinalPrice
	existing.DurationPrice = req.DurationPrice
	existing.TokenAmount = req.TokenAmount
	existing.BonusToken = req.BonusToken
	existing.PointsActive = req.PointsActive
	existing.PointsValue = req.PointsValue
	existing.BonusPoints = req.BonusPoints
	existing.IsActive = req.IsActive
	existing.UpdatedAt = time.Now().UTC()

	result, err := s.durationRepository.Update(ctx, nil, existing)
	if err != nil { return dto_response.GetPlanDurationResponse{}, myerror.DatabaseError(err) }
	count, _ := s.durationRepository.CountMappings(ctx, nil, result.ID)
	return toDurationResponse(result, int(count)), nil
}

func (s *planDurationService) Delete(ctx context.Context, planID, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil { return myerror.InvalidRequest(err) }
	existing, err := s.durationRepository.GetByID(ctx, nil, parsedID)
	if err != nil { return myerror.RecordNotFound("plan duration") }
	if existing.PlanID.String() != planID { return myerror.RecordNotFound("plan duration") }

	memberCount, _ := s.durationRepository.CountActiveMembers(ctx, nil, parsedID)
	if memberCount > 0 { return myerror.New(fmt.Sprintf("has %d active members", memberCount), myerror.Error_InvalidRequest) }
	mappingCount, _ := s.durationRepository.CountMappings(ctx, nil, parsedID)
	if mappingCount > 0 { return myerror.New(fmt.Sprintf("has %d mappings", mappingCount), myerror.Error_InvalidRequest) }

	return s.durationRepository.Delete(ctx, nil, parsedID)
}

func toDurationResponse(d entity.PlanDuration, mappingCount int) dto_response.GetPlanDurationResponse {
	createdAt, updatedAt := "", ""
	if !d.CreatedAt.IsZero() { createdAt = d.CreatedAt.Format(time.RFC3339) }
	if !d.UpdatedAt.IsZero() { updatedAt = d.UpdatedAt.Format(time.RFC3339) }
	return dto_response.GetPlanDurationResponse{
		ID: d.ID.String(), PlanID: d.PlanID.String(), DurationMonths: d.DurationMonths,
		Price: d.Price, DiscountPct: d.DiscountPct, FinalPrice: d.FinalPrice, DurationPrice: d.DurationPrice,
		TokenAmount: d.TokenAmount, BonusToken: d.BonusToken, PointsActive: d.PointsActive, PointsValue: d.PointsValue, BonusPoints: d.BonusPoints,
		IsActive: d.IsActive, MappingCount: mappingCount, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}
