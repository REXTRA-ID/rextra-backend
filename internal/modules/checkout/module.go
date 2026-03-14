package checkout

import (
	"context"

	"rextra-backend/internal/entity"
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/checkout/controller"
	"rextra-backend/internal/modules/checkout/repository"
	"rextra-backend/internal/modules/checkout/routes"
	"rextra-backend/internal/modules/checkout/service"
	"rextra-backend/internal/pkg/tripay"

	membershipRepo "rextra-backend/internal/modules/membership/repository"
	promoService "rextra-backend/internal/modules/promo/service"
	userRepo "rextra-backend/internal/modules/user/repository"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	dto_request "rextra-backend/internal/dto/request"
)

// Adapters
type planRepoAdapter struct {
	repo membershipRepo.MembershipPlanRepository
}
func (a *planRepoAdapter) GetByIDWithDurations(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.MembershipPlans, error) {
	return a.repo.GetByID(ctx, tx, id)
}

type promoRepoAdapter struct {
	svc promoService.DiscountService
}
func (a *promoRepoAdapter) CountEligibleVouchers(ctx context.Context, planID string) (int, error) {
	// For now return total active discounts
	return 0, nil
}
func (a *promoRepoAdapter) ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, durationID string, subtotal int64) (int64, error) {
	res, err := a.svc.ValidateCode(ctx, dto_request.ValidateDiscountRequest{Code: code, PlanID: planID, DurationID: durationID}, "")
	if err != nil { return 0, err }
	return res.DiscountAmount, nil
}
func (a *promoRepoAdapter) RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID uuid.UUID) error {
	// Handled by payment service after success payment callback
	return nil
}

type userRepoAdapter struct {
	repo userRepo.UserRepository
}
func (a *userRepoAdapter) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entity.User, error) {
	return a.repo.GetById(ctx, tx, id.String())
}

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware, tripayClient *tripay.TripayClient, promoSvc promoService.DiscountService) {
	repo := repository.NewCheckoutRepository(db)
	
	planAdapter := &planRepoAdapter{repo: membershipRepo.NewMembershipPlanRepository(db)}
	promoAdapter := &promoRepoAdapter{svc: promoSvc}
	userAdapter := &userRepoAdapter{repo: userRepo.NewUserRepository(db)}

	svc := service.NewCheckoutService(repo, planAdapter, promoAdapter, userAdapter, tripayClient, db)
	ctrl := controller.NewCheckoutController(svc)

	routes.ServeCheckout(server, ctrl, mw)
}
