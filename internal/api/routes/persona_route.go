package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServePersona(app *gin.Engine, personacontroller controller.PersonaController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/persona")
	{
		routes.POST("/create", middleware.Authenticate(), personacontroller.Create)
		routes.GET("", middleware.Authenticate(), personacontroller.Get)
		routes.PUT("/mission/complete", middleware.Authenticate(), personacontroller.UpdateMission)
	}
}
