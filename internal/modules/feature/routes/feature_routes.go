package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/feature/controller"

	"github.com/gin-gonic/gin"
)

func ServeFeature(app *gin.Engine, featureController controller.FeatureController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/feature")
	{
		routes.GET("", middleware.Authenticate(), featureController.GetAll)
		routes.GET("/:featureId", middleware.Authenticate(), featureController.GetById)
		routes.POST("", middleware.Authenticate(), featureController.Create)
	}
}
