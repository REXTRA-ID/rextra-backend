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
	entitlementRepo "rextra-backend/internal/modules/user_entitlement/repository"
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
		GetCatalog(ctx context.Context, userID string) (dto_response.MembershipCatalogResponse, error)
		GetStarterPlan(ctx context.Context) (entity.MembershipPlans, error)
	}

	membershipPlanService struct {
		planRepository  repository.MembershipPlanRepository
		membershipRepo  repository.UserMembershipRepository
		entitlementRepo entitlementRepo.UserEntitlementRepository
	}
)

func NewMembershipPlanService(
	planRepo repository.MembershipPlanRepository,
	membershipRepo repository.UserMembershipRepository,
	entitlementRepo entitlementRepo.UserEntitlementRepository,
) MembershipPlanService {
	return &membershipPlanService{
		planRepository:  planRepo,
		membershipRepo:  membershipRepo,
		entitlementRepo: entitlementRepo,
	}
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
	newPlan.MarketingIntro = req.MarketingIntro
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
	existing.MarketingIntro = req.MarketingIntro
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

func (s *membershipPlanService) GetCatalog(ctx context.Context, userID string) (dto_response.MembershipCatalogResponse, error) {
	var membership entity.Memberships
	if userID != "" {
		parsedUserID, err := uuid.Parse(userID)
		if err == nil {
			membership, _ = s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
		}
	}

	allPlans, err := s.planRepository.GetAllWithDurations(ctx, nil)
	if err != nil { return dto_response.MembershipCatalogResponse{}, myerror.DatabaseError(err) }

	tierRank := map[entity.PlanName]int{
		entity.PlanNameStandard: 0,
		entity.PlanNameStarter:  0,
		entity.PlanNameBasic:    1,
		entity.PlanNamePro:      2,
		entity.PlanNameMax:      3,
	}

	currentRank := tierRank[membership.PlanName]
	contextMsg := ""
	if membership.IsActive && membership.PlanName != entity.PlanNameStandard && membership.PlanName != entity.PlanNameStarter {
		dur := 0; if membership.DurationMonths != nil { dur = *membership.DurationMonths }
		contextMsg = fmt.Sprintf("Kamu sedang berlangganan %s — %d bulan.", membership.PlanName, dur)
	}

	resp := dto_response.MembershipCatalogResponse{}
	if contextMsg != "" {
		planIDStr := ""; if membership.PlanID != nil { planIDStr = membership.PlanID.String() }
		resp.CurrentContext = &dto_response.CurrentSubscriptionContext{
			Message: contextMsg, CurrentPlanID: planIDStr, CurrentPlanName: string(membership.PlanName),
		}
	}

	for _, p := range allPlans {
		if p.Status != entity.PlanStatusActive || p.PlanName == entity.PlanNameStandard || p.PlanName == entity.PlanNameStarter { continue }
		
		targetRank := tierRank[p.PlanName]
		ctaLabel := "Pilih " + string(p.PlanName) + " Plan"
		changeType := "PEMBELIAN_BARU"
		isCurrent := membership.PlanName == p.PlanName && membership.IsActive

		if membership.IsActive && membership.PlanName != entity.PlanNameStandard && membership.PlanName != entity.PlanNameStarter {
			if isCurrent {
				ctaLabel = "Tambah Durasi"
				changeType = "RENEWAL"
			} else if targetRank > currentRank {
				ctaLabel = "Upgrade ke plan ini"
				changeType = "UPGRADE"
			} else {
				ctaLabel = "Downgrade ke plan ini"
				changeType = "DOWNGRADE"
			}
		}

		item := dto_response.PlanCatalogItem{
			ID: p.ID.String(), PlanName: string(p.PlanName), TierLabel: p.TierLabel, EmblemKey: p.EmblemKey,
			MarketingIntro: p.MarketingIntro, PricingInfo: fmt.Sprintf("Mulai dari Rp%d", p.BasePrice1M),
			TokenBonusInfo: fmt.Sprintf("%d Token", p.BaseToken1M), ThemeColor: s.getThemeColor(p.PlanName),
			CTA: dto_response.PlanCTAInfo{Label: ctaLabel, ChangeType: changeType, IsCurrent: isCurrent},
		}
		
		if p.PlanName == entity.PlanNameBasic { item.Label = "Paling Direkomendasikan" }
		if p.PlanName == entity.PlanNamePro { item.Label = "Lebih Lengkap" }
		if p.PlanName == entity.PlanNameMax { item.Label = "Akses Tanpa Batas" }

		if len(p.PlanDurations) > 0 {
			firstDurID := p.PlanDurations[0].ID
			mappings, _ := s.entitlementRepo.GetEntitlementsByPlanDurationID(ctx, &firstDurID)
			
			benefitGroups := make(map[uuid.UUID]*dto_response.FeatureBenefitGroup)
			var groupOrder []uuid.UUID
			for _, m := range mappings {
				e := m.Entitlement
				f := e.Feature
				if _, ok := benefitGroups[f.ID]; !ok {
					benefitGroups[f.ID] = &dto_response.FeatureBenefitGroup{FeatureName: f.Name, Description: f.Description, IconKey: f.Slug}
					groupOrder = append(groupOrder, f.ID)
				}
				group := benefitGroups[f.ID]
				if f.Type == entity.FeatureTypeSingle {
					group.AccessLabel = s.getAccessLabel(e, entity.UserEntitlementQuota{})
				} else {
					group.AccessLabel = "Bervariasi"
					group.SubFeatures = append(group.SubFeatures, dto_response.SubFeatureBenefit{Name: e.Name, AccessLabel: s.getAccessLabel(e, entity.UserEntitlementQuota{})})
				}
			}
			for _, id := range groupOrder { item.Benefits = append(item.Benefits, *benefitGroups[id]) }
		}
		resp.Plans = append(resp.Plans, item)
	}
	return resp, nil
}

func (s *membershipPlanService) getThemeColor(plan entity.PlanName) string {
	switch plan {
	case entity.PlanNameStarter: return "#E3F2FD"
	case entity.PlanNameBasic: return "#E8F5E9"
	case entity.PlanNamePro: return "#F3E5F5"
	case entity.PlanNameMax: return "#FFF8E1"
	default: return "#F5F5F5"
	}
}

func (s *membershipPlanService) getAccessLabel(e entity.Entitlement, q entity.UserEntitlementQuota) string {
	switch e.RestrictionType {
	case entity.RestrictionUnlimited: return "Tanpa Batas"
	case entity.RestrictionTokenGated: return "Berbasis Token"
	case entity.RestrictionFrequencyLimited: 
		if q.ID != uuid.Nil { return fmt.Sprintf("%d Sisa Kuota", q.QuotaRemaining) }
		return "Kuota Terbatas"
	case entity.RestrictionLocked: return "Terkunci"
	default: return "Terbatas"
	}
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
