package membership

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/membership/controller"
	"rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/membership/routes"
	"rextra-backend/internal/modules/membership/service"

	entitlementRepo "rextra-backend/internal/modules/entitlement/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	planRepo := repository.NewMembershipPlanRepository(db)
	durationRepo := repository.NewPlanDurationRepository(db)
	mappingRepo := repository.NewDurationAccessMappingRepository(db)
	entRepo := entitlementRepo.NewEntitlementRepository(db)

	planService := service.NewMembershipPlanService(planRepo)
	durationService := service.NewPlanDurationService(planRepo, durationRepo)
	mappingService := service.NewDurationAccessMappingService(mappingRepo, durationRepo, entRepo)

	planController := controller.NewMembershipPlanController(planService)
	durationController := controller.NewPlanDurationController(durationService)
	mappingController := controller.NewDurationAccessMappingController(mappingService)

	routes.ServeMembership(server, planController, durationController, mappingController, middleware)
}
