package userRoutes

import (
	userController "rextra-backend/internal/api/user/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(app *gin.Engine, usercontroller userController.UserController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/users")
	{
		routes.GET("/:id", usercontroller.GetById)
	}
}
