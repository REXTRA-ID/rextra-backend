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
	// Repositories
	planRepository := repository.NewMembershipPlanRepository(db)
	durationRepository := repository.NewMembershipDurationRepository(db)
	mappingRepository := repository.NewDurationAccessMappingRepository(db)
	cycleRepository := repository.NewSubscriptionCycleRepository(db)

	// entitlement repository dari modul entitlement — dibutuhkan mapping service
	// untuk load snapshot data entitlement saat Create mapping.
	// Go structural typing memastikan implementasi ini memenuhi interface lokal
	// repository.EntitlementRepository di modul membership.
	entitlementRepository := entitlementRepo.NewEntitlementRepository(db)

	// Services
	planService := service.NewMembershipPlanService(planRepository)
	durationService := service.NewMembershipDurationService(durationRepository)
	mappingService := service.NewDurationAccessMappingService(mappingRepository, durationRepository, entitlementRepository)
	cycleService := service.NewSubscriptionCycleService(cycleRepository)

	// Controllers
	planController := controller.NewMembershipPlanController(planService)
	durationController := controller.NewPlanDurationController(durationService)
	mappingController := controller.NewDurationAccessMappingController(mappingService)
	cycleController := controller.NewSubscriptionCycleController(cycleService)

	// Routes
	routes.ServeMembership(server, planController, durationController, mappingController, cycleController, middleware)
}
