package entitlement

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/entitlement/controller"
	"rextra-backend/internal/modules/entitlement/repository"
	"rextra-backend/internal/modules/entitlement/routes"
	"rextra-backend/internal/modules/entitlement/service"

	actionCategoryRepo "rextra-backend/internal/modules/action_category/repository"
	featureRepo "rextra-backend/internal/modules/feature/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware) {
	entitlementRepository := repository.NewEntitlementRepository(db)

	featureRepository := featureRepo.NewFeatureRepository(db)
	subFeatureRepository := featureRepo.NewSubFeatureRepository(db)
	actionCategoryRepository := actionCategoryRepo.NewActionCategoryRepository(db)

	entitlementService := service.NewEntitlementService(
		entitlementRepository,
		featureRepository,
		subFeatureRepository,
		actionCategoryRepository,
	)

	entitlementController := controller.NewEntitlementController(entitlementService)

	routes.ServeEntitlement(server, entitlementController, mw)
}
