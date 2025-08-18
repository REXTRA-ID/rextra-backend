package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeRiasec(app *gin.Engine, riaseccontroller controller.RiasecController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/user/:id/riasec")
	{
		routes.POST("", middleware.Authenticate(), riaseccontroller.Create)
		routes.GET("", middleware.Authenticate(), riaseccontroller.GetByUserID)
		routes.PUT("", middleware.Authenticate(), riaseccontroller.Update)
	}
}
