package feature

import (
	"rextra-backend/internal/middleware"
	fc "rextra-backend/internal/modules/feature/controller"
	fr "rextra-backend/internal/modules/feature/repository"
	"rextra-backend/internal/modules/feature/routes"
	fs "rextra-backend/internal/modules/feature/service"
	sfc "rextra-backend/internal/modules/sub_feature/controller"
	sfr "rextra-backend/internal/modules/sub_feature/repository"
	sfs "rextra-backend/internal/modules/sub_feature/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repositories
	featureRepository := fr.NewFeatureRepository(db)
	// NewSubFeatureRepository berada di package feature/repository (bukan modul terpisah)
	// karena SubFeature adalah sub-resource dari Feature dalam satu modul yang sama.
	subFeatureRepository := sfr.NewSubFeatureRepository(db)

	// Services
	featureService := fs.NewFeatureService(featureRepository)
	subFeatureService := sfs.NewSubFeatureService(featureRepository, subFeatureRepository)

	// Controllers
	featureController := fc.NewFeatureController(featureService)
	subFeatureController := sfc.NewSubFeatureController(subFeatureService)

	// Routes
	routes.ServeFeature(server, featureController, subFeatureController, middleware)
}
