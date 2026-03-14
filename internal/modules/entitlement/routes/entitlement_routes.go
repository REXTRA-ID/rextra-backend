package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/entitlement/controller"

	"github.com/gin-gonic/gin"
)

func ServeEntitlement(
	app *gin.Engine,
	ctrl controller.EntitlementController,
	mw middleware.Middleware,
) {
	routes := app.Group("/api/v1/entitlement")
	routes.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		routes.GET("", ctrl.GetAll)
		routes.POST("", ctrl.Create)
		routes.GET("/:entitlementId", ctrl.GetById)
		routes.PATCH("/:entitlementId/restriction", ctrl.UpdateRestriction)
		routes.DELETE("/:entitlementId", ctrl.Delete)
	}
}
