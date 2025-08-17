package authRoutes

import (
	authController "rextra-backend/internal/api/auth/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(app *gin.Engine, authcontroller authController.AuthController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/auth")
	{
		routes.POST("/login", authcontroller.Login)
		routes.POST("/register", authcontroller.Register)
		routes.GET("/verify", authcontroller.Verify)
		routes.GET("/me", middleware.Authenticate(), authcontroller.Me)

		routes.POST("/google", middleware.Authenticate(), authcontroller.LoginWithGoogle)
		// routes.GET("/google/callback", authcontroller.CallbackGoogle)
	}
}
