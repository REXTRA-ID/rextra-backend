package entitlement

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/action_category/repository"
	"rextra-backend/internal/modules/entitlement/controller"
	erepo "rextra-backend/internal/modules/entitlement/repository"
	"rextra-backend/internal/modules/entitlement/routes"
	"rextra-backend/internal/modules/entitlement/service"
	frepo "rextra-backend/internal/modules/feature/repository"
	sfrepo "rextra-backend/internal/modules/sub_feature/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repositories
	entitlementRepository := erepo.NewEntitlementRepository(db)
	featureRepository := frepo.NewFeatureRepository(db)
	subFeatureRepository := sfrepo.NewSubFeatureRepository(db)
	actionCategoryRepository := repository.NewActionCategoryRepository(db)

	// Service
	entitlementService := service.NewEntitlementService(
		entitlementRepository,
		featureRepository,
		subFeatureRepository,
		actionCategoryRepository,
	)

	// Controller
	entitlementController := controller.NewEntitlementController(entitlementService)

	// Routes
	routes.ServeEntitlement(server, entitlementController, middleware)
}
