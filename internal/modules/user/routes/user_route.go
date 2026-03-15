package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth/controller"

	"github.com/gin-gonic/gin"
)

func ServeAuth(app *gin.Engine, authcontroller controller.AuthController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/users")
	{
		routes.GET("/:id", authcontroller.Me)
	}
}
