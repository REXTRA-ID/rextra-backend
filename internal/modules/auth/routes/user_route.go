package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth/controller"

	"github.com/gin-gonic/gin"
)

func ServeUser(app *gin.Engine, usercontroller controller.UserController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/users")
	{
		routes.GET("/:id", usercontroller.GetById)
	}
}
