package membership

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/membership/controller"
	"rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/membership/routes"
	"rextra-backend/internal/modules/membership/service"

	entitlementRepo "rextra-backend/internal/modules/entitlement/repository"
	userEntitlementRepo "rextra-backend/internal/modules/user_entitlement/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware) {
	planRepo := repository.NewMembershipPlanRepository(db)
	durationRepo := repository.NewPlanDurationRepository(db)
	membershipRepo := repository.NewUserMembershipRepository(db)
	userEntitlementRepository := userEntitlementRepo.NewUserEntitlementRepository(db)
	entRepo := entitlementRepo.NewEntitlementRepository(db)

	planSvc := service.NewMembershipPlanService(planRepo, membershipRepo, userEntitlementRepository)
	durationSvc := service.NewPlanDurationService(planRepo, durationRepo)
	mappingSvc := service.NewDurationAccessMappingService(repository.NewDurationAccessMappingRepository(db), durationRepo, entRepo)

	planCtrl := controller.NewMembershipPlan(planSvc)
	durationCtrl := controller.NewPlanDurationController(durationSvc)
	mappingCtrl := controller.NewDurationAccessMappingController(mappingSvc)

	routes.ServeMembership(server, planCtrl, durationCtrl, mappingCtrl, mw)
}
