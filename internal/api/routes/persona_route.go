package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServePersona(app *gin.Engine, personacontroller controller.PersonaController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/user/:id/persona")
	{
		routes.POST("", middleware.Authenticate(), personacontroller.Create)
		routes.GET("", middleware.Authenticate(), personacontroller.GetByUserID)
		routes.PUT("", middleware.Authenticate(), personacontroller.Update)
	}
}
