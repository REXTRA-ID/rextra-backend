package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/feature/controller"

	"github.com/gin-gonic/gin"
)

func ServeFeature(
	app *gin.Engine,
	featureController controller.FeatureController,
	subFeatureController controller.SubFeatureController,
	middleware middleware.Middleware,
) {
	feature := app.Group("/api/v1/feature")
	feature.Use(middleware.Authenticate())
	{
		feature.GET("", featureController.GetAll)
		feature.GET("/:featureId", featureController.GetById)
		feature.POST("", featureController.Create)
		feature.PUT("/:featureId", featureController.Update)
		feature.DELETE("/:featureId", featureController.Delete)

		subFeature := feature.Group("/:featureId/sub-feature")
		{
			subFeature.GET("", subFeatureController.GetAll)
			subFeature.GET("/:subFeatureId", subFeatureController.GetById)
			subFeature.POST("", subFeatureController.Create)
			subFeature.PUT("/:subFeatureId", subFeatureController.Update)
			subFeature.DELETE("/:subFeatureId", subFeatureController.Delete)
		}
	}
}
