package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	membershipRepo "rextra-backend/internal/modules/membership/repository"
	userRepo "rextra-backend/internal/modules/user/repository"
	entitlementRepo "rextra-backend/internal/modules/user_entitlement/repository"
	tokenRepo "rextra-backend/internal/modules/token/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	UserMembershipService interface {
		GetAll(ctx context.Context, req dto_request.UserMembershipFilterRequest) (dto_response.PaginatedUserMembershipResponse, error)
		GetById(ctx context.Context, membershipID string) (dto_response.GetUserMembershipResponse, error)
		GetMyMembership(ctx context.Context, userID string) (dto_response.GetUserMembershipResponse, error)
		ClaimStarter(ctx context.Context, userID string) error
		GetDashboard(ctx context.Context, userID string) (dto_response.MembershipDashboardResponse, error)
	}

	userMembershipService struct {
		membershipRepo     membershipRepo.UserMembershipRepository
		membershipPlanRepo membershipRepo.MembershipPlanRepository
		userRepo           userRepo.UserRepository
		entitlementRepo    entitlementRepo.UserEntitlementRepository
		tokenRepo          tokenRepo.TokenWalletRepository
	}
)

func NewUserMembershipService(
	membershipRepo membershipRepo.UserMembershipRepository,
	membershipPlanRepo membershipRepo.MembershipPlanRepository,
	userRepo userRepo.UserRepository,
	entitlementRepo entitlementRepo.UserEntitlementRepository,
	tokenRepo tokenRepo.TokenWalletRepository,
) UserMembershipService {
	return &userMembershipService{
		membershipRepo:     membershipRepo,
		membershipPlanRepo: membershipPlanRepo,
		userRepo:           userRepo,
		entitlementRepo:    entitlementRepo,
		tokenRepo:          tokenRepo,
	}
}

func (s *userMembershipService) GetAll(ctx context.Context, req dto_request.UserMembershipFilterRequest) (dto_response.PaginatedUserMembershipResponse, error) {
	page := req.Page; if page == 0 { page = 1 }
	pageSize := req.PageSize; if pageSize == 0 { pageSize = 10 }
	offset := (page - 1) * pageSize

	models, total, err := s.membershipRepo.GetAllPaginated(ctx, nil, offset, pageSize, req.PlanName, req.Search, req.IsActive)
	if err != nil { return dto_response.PaginatedUserMembershipResponse{}, myerror.DatabaseError(err) }

	var data []dto_response.GetUserMembershipListResponse
	for _, m := range models { data = append(data, toListResponse(m)) }

	return dto_response.PaginatedUserMembershipResponse{
		Data: data, Total: int(total), Page: page, PageSize: pageSize, TotalPages: int(math.Ceil(float64(total) / float64(pageSize))),
	}, nil
}

func (s *userMembershipService) GetById(ctx context.Context, membershipID string) (dto_response.GetUserMembershipResponse, error) {
	parsedID, err := uuid.Parse(membershipID)
	if err != nil { return dto_response.GetUserMembershipResponse{}, myerror.InvalidRequest(err) }
	membership, err := s.membershipRepo.GetByID(ctx, nil, parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetUserMembershipResponse{}, myerror.RecordNotFound("membership") }
		return dto_response.GetUserMembershipResponse{}, myerror.DatabaseError(err)
	}
	return toDetailResponse(membership), nil
}

func (s *userMembershipService) GetMyMembership(ctx context.Context, userID string) (dto_response.GetUserMembershipResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return dto_response.GetUserMembershipResponse{}, myerror.InvalidRequest(err) }
	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return dto_response.GetUserMembershipResponse{}, myerror.RecordNotFound("membership") }
		return dto_response.GetUserMembershipResponse{}, myerror.DatabaseError(err)
	}
	return toDetailResponse(membership), nil
}

func (s *userMembershipService) ClaimStarter(ctx context.Context, userID string) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return myerror.InvalidRequest(err) }

	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	if err != nil { return myerror.DatabaseError(err) }

	if membership.HasClaimedStarter {
		return myerror.New("anda sudah pernah mengklaim starter plan", myerror.Error_InvalidRequest)
	}

	starterPlan, err := s.membershipPlanRepo.GetByPlanName(ctx, nil, entity.PlanName(entity.PLANSTARTER))
	if err != nil { return myerror.RecordNotFound("starter plan") }

	now := time.Now().UTC()
	expiredAt := now.AddDate(0, starterPlan.StarterDurationMonths, 0)

	membership.PlanID = &starterPlan.ID
	membership.PlanName = starterPlan.PlanName
	
	// Find the 1-month duration for Starter plan to link benefits
	starterDurations, err := s.membershipPlanRepo.GetDurationsByPlanID(ctx, nil, starterPlan.ID)
	if err == nil {
		for _, d := range starterDurations {
			if d.DurationMonths == 1 {
				membership.DurationID = &d.ID
				mDur := 1
				membership.DurationMonths = &mDur
				break
			}
		}
	}

	membership.StartedAt = &now
	membership.ExpiredAt = &expiredAt
	membership.IsActive = true
	membership.HasClaimedStarter = true
	membership.UpdatedAt = now

	_, err = s.membershipRepo.Update(ctx, nil, membership)
	return err
}

func (s *userMembershipService) GetDashboard(ctx context.Context, userID string) (dto_response.MembershipDashboardResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil { return dto_response.MembershipDashboardResponse{}, myerror.InvalidRequest(err) }

	// 1. Get Active Membership
	membership, err := s.membershipRepo.GetByUserID(ctx, nil, parsedUserID)
	if err != nil { return dto_response.MembershipDashboardResponse{}, myerror.DatabaseError(err) }

	// 2. Get Plan Details (for visual assets)
	var tierLabel, emblemKey string
	if membership.PlanID != nil {
		plan, err := s.membershipPlanRepo.GetByID(ctx, nil, *membership.PlanID)
		if err == nil {
			tierLabel = plan.TierLabel
			emblemKey = plan.EmblemKey
		}
	}

	// 3. Get Real Token Balance
	wallet, _ := s.tokenRepo.GetOrCreateByUserID(ctx, nil, parsedUserID)

	// 4. Construct Dashboard Response
	remainingDays := membership.RemainingDays()
	isNearExpiry := membership.IsActive && remainingDays <= 7 && membership.PlanName != entity.PlanName(entity.PLANSTARTER)
	
	activeUntil := ""
	if membership.ExpiredAt != nil { activeUntil = membership.ExpiredAt.Format("02 January 2006") }

	statusLabel := "REXTRA CLUB"
	if membership.PlanName == entity.PlanName(entity.PLANSTARTER) || membership.PlanName == entity.PlanName(entity.PLANSTANDARD) {
		statusLabel = "REXTRA NON CLUB"
	}

	resp := dto_response.MembershipDashboardResponse{
		Membership: dto_response.MembershipStatusInfo{
			PlanName: string(membership.PlanName), TierLabel: tierLabel, EmblemKey: emblemKey, StatusLabel: statusLabel, ActiveUntil: &activeUntil, RemainingDays: remainingDays,
		},
		Wallet: dto_response.TokenWalletInfo{TokenBalance: wallet.Balance},
		UIState: dto_response.MembershipUIState{
			IsNearExpiry: isNearExpiry, IsExpired: !membership.IsActive, HasClaimedStarter: membership.HasClaimedStarter, ThemeColor: getThemeColor(membership.PlanName),
		},
	}

	// 5. Build Benefit Tree
	if membership.DurationID != nil {
		mappings, _ := s.entitlementRepo.GetEntitlementsByPlanDurationID(ctx, membership.DurationID)
		quotas, _ := s.entitlementRepo.GetQuotasByMembershipID(ctx, membership.ID)
		quotaMap := make(map[string]entity.UserEntitlementQuota)
		for _, q := range quotas { quotaMap[q.EntitlementKey] = q }

		// Grouping
		benefitGroups := make(map[uuid.UUID]*dto_response.FeatureBenefitGroup)
		var groupOrder []uuid.UUID

		for _, m := range mappings {
			e := m.Entitlement
			f := e.Feature
			
			if _, ok := benefitGroups[f.ID]; !ok {
				benefitGroups[f.ID] = &dto_response.FeatureBenefitGroup{
					FeatureName: f.Name, Description: f.Description, IconKey: f.Slug,
				}
				groupOrder = append(groupOrder, f.ID)
			}
			
			group := benefitGroups[f.ID]
			
			if f.Type == entity.FeatureTypeSingle {
				group.AccessType = string(e.RestrictionType)
				group.AccessLabel = getAccessLabel(e, quotaMap[e.Key])
			} else {
				group.AccessType = "varied"
				group.AccessLabel = "Bervariasi"
				
				sub := dto_response.SubFeatureBenefit{Name: e.Name, AccessType: string(e.RestrictionType), AccessLabel: getAccessLabel(e, quotaMap[e.Key])}
				if q, ok := quotaMap[e.Key]; ok && e.RestrictionType == entity.RestrictionFrequencyLimited {
					left := q.QuotaRemaining; sub.QuotaLeft = &left
				}
				group.SubFeatures = append(group.SubFeatures, sub)
			}
		}

		for _, id := range groupOrder { resp.Benefits = append(resp.Benefits, *benefitGroups[id]) }
	}

	return resp, nil
}

func getThemeColor(plan entity.PlanName) string {
	switch plan {
	case entity.PlanName(entity.PLANSTARTER): return "#E3F2FD" // Light Blue
	case entity.PlanName(entity.PLANBASIC): return "#E8F5E9"   // Green
	case entity.PlanName(entity.PLANPRO): return "#F3E5F5"     // Purple
	case entity.PlanName(entity.PLANMAX): return "#FFF8E1"     // Gold/Amber
	default: return "#F5F5F5" // Grey
	}
}

func getAccessLabel(e entity.Entitlement, q entity.UserEntitlementQuota) string {
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

func toDetailResponse(m entity.Memberships) dto_response.GetUserMembershipResponse {
	resp := dto_response.GetUserMembershipResponse{
		MembershipID: m.ID.String(), UserID: m.UserID.String(), PlanName: string(m.PlanName), IsActive: m.IsActive, AutoRenew: m.AutoRenew,
		RemainingDays: m.RemainingDays(), PaidCycleCount: m.PaidCycleCount, CurrentTokenBalance: m.CurrentTokenBalance, CurrentPoinBalance: m.CurrentPoinBalance, EntitlementCount: m.EntitlementCount, UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
	if m.PlanID != nil { resp.PlanID = m.PlanID.String() }
	if m.DurationID != nil { resp.DurationID = m.DurationID.String() }
	if m.DurationMonths != nil { resp.DurationMonths = m.DurationMonths }
	if m.StartedAt != nil { s := m.StartedAt.Format(time.RFC3339); resp.StartedAt = &s }
	if m.ExpiredAt != nil { e := m.ExpiredAt.Format(time.RFC3339); resp.ExpiredAt = &e }
	return resp
}

func toListResponse(m entity.Memberships) dto_response.GetUserMembershipListResponse {
	resp := dto_response.GetUserMembershipListResponse{
		MembershipID: m.ID.String(), UserID: m.UserID.String(), PlanName: string(m.PlanName), IsActive: m.IsActive,
		RemainingDays: m.RemainingDays(), PaidCycleCount: m.PaidCycleCount, UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
	if m.DurationMonths != nil { resp.DurationMonths = m.DurationMonths }
	if m.ExpiredAt != nil { e := m.ExpiredAt.Format(time.RFC3339); resp.ExpiredAt = &e }
	return resp
}
