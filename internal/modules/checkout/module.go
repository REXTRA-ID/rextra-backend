package checkout

import (
	"context"
	"encoding/json"
	"strings"

	"rextra-backend/internal/entity"
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/checkout/controller"
	"rextra-backend/internal/modules/checkout/repository"
	"rextra-backend/internal/modules/checkout/routes"
	"rextra-backend/internal/modules/checkout/service"
	"rextra-backend/internal/pkg/tripay"

	membershipRepo "rextra-backend/internal/modules/membership/repository"
	promoRepo "rextra-backend/internal/modules/promo/repository"
	userRepo "rextra-backend/internal/modules/user/repository"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type planRepoAdapter struct {
	repo membershipRepo.MembershipPlanRepository
}
func (a *planRepoAdapter) GetByIDWithDurations(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error) {
	return a.repo.GetByID(ctx, tx, id)
}

type promoRepoAdapter struct {
	discountRepo   promoRepo.DiscountRepository
	redemptionRepo promoRepo.DiscountRedemptionRepository
	planRepo       membershipRepo.MembershipPlanRepository
}
func (a *promoRepoAdapter) CountEligibleVouchers(ctx context.Context, planID string) (int, error) {
	return 0, nil // Simple placeholder
}
func (a *promoRepoAdapter) ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, subtotal int64) (int64, error) {
	discount, err := a.discountRepo.GetByCode(ctx, nil, code)
	if err != nil { return 0, err }
	if !discount.IsValid() { return 0, gorm.ErrRecordNotFound }
	
	if discount.MembershipPlanTargets != nil {
		plan, err := a.planRepo.GetByID(ctx, nil, uuid.MustParse(planID))
		if err == nil {
			var targets []string
			json.Unmarshal(discount.MembershipPlanTargets, &targets)
			matched := false
			for _, t := range targets { if strings.EqualFold(t, string(plan.PlanName)) { matched = true; break } }
			if !matched { return 0, gorm.ErrRecordNotFound }
		}
	}
	return discount.CalculateDiscount(subtotal), nil
}
func (a *promoRepoAdapter) RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID string) error {
	discount, err := a.discountRepo.GetByCode(ctx, nil, code)
	if err != nil { return err }
	
	redemption := entity.DiscountRedemption{
		DiscountID: discount.ID, UserID: userID, TransactionID: transactionID, CodeSnapshot: discount.Code, Status: entity.RedemptionStatusApplied,
	}
	_, err = a.redemptionRepo.Create(ctx, nil, redemption)
	if err != nil { return err }
	return a.discountRepo.IncrementRedemption(ctx, nil, discount.ID)
}

type userRepoAdapter struct {
	repo userRepo.UserRepository
}
func (a *userRepoAdapter) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.User, error) {
	return a.repo.GetById(ctx, tx, id.String())
}

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware, tripayClient *tripay.TripayClient) {
	repo := repository.NewCheckoutRepository(db)
	
	planRepo := membershipRepo.NewMembershipPlanRepository(db)
	planAdapter := &planRepoAdapter{repo: planRepo}
	
	promoAdapter := &promoRepoAdapter{
		discountRepo: promoRepo.NewDiscountRepository(db),
		redemptionRepo: promoRepo.NewDiscountRedemptionRepository(db),
		planRepo: planRepo,
	}
	
	userAdapter := &userRepoAdapter{repo: userRepo.NewUserRepository(db)}

	svc := service.NewCheckoutService(repo, planAdapter, promoAdapter, userAdapter, tripayClient, db)
	ctrl := controller.NewCheckoutController(svc)

	routes.ServeCheckout(server, ctrl, mw)
}
