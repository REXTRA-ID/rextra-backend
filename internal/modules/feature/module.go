package feature

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/feature/controller"
	"rextra-backend/internal/modules/feature/repository"
	"rextra-backend/internal/modules/feature/routes"
	"rextra-backend/internal/modules/feature/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	featureRepository := repository.NewFeatureRepository(db)

	featureService := service.NewFeatureService(featureRepository)

	featureController := controller.NewFeatureController(featureService)

	routes.ServeFeature(server, featureController, middleware)
}
