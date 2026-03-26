package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/entitlement/controller"

	"github.com/gin-gonic/gin"
)

func ServeEntitlement(app *gin.Engine, entitlementController controller.EntitlementController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/entitlement")
	routes.Use(middleware.Authenticate())
	{
		routes.GET("", entitlementController.GetAll)
		routes.GET("/:entitlementId", entitlementController.GetById)
		routes.POST("", entitlementController.Create)
		routes.DELETE("/:entitlementId", entitlementController.Delete)
	}
}
