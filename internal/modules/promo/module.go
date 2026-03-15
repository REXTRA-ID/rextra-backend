package promo

import (
	"context"

	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/promo/controller"
	"rextra-backend/internal/modules/promo/repository"
	"rextra-backend/internal/modules/promo/routes"
	"rextra-backend/internal/modules/promo/service"

	membershipRepo "rextra-backend/internal/modules/membership/repository"

	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type planDurationAdapter struct {
	repo membershipRepo.PlanDurationRepository
}

func (a *planDurationAdapter) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (service.PlanDurationInfo, error) {
	duration, err := a.repo.GetByID(ctx, tx, id)
	if err != nil { return service.PlanDurationInfo{}, err }
	planName := ""
	if duration.Plan != nil { planName = string(duration.Plan.PlanName) }
	return service.PlanDurationInfo{ID: duration.ID, PlanID: duration.PlanID, PlanName: planName, FinalPrice: duration.FinalPrice, DurationMonths: duration.DurationMonths}, nil
}

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware) {
	discountRepository := repository.NewDiscountRepository(db)
	redemptionRepository := repository.NewDiscountRedemptionRepository(db)
	planDurationRepo := membershipRepo.NewPlanDurationRepository(db)
	planDurationReaderAdapter := &planDurationAdapter{repo: planDurationRepo}

	discountService := service.NewDiscountService(discountRepository, redemptionRepository, planDurationReaderAdapter, db)
	discountController := controller.NewDiscountController(discountService)

	routes.ServePromo(server, discountController, mw)
}
