package routes

import (
	"rextra-backend/internal/modules/auth/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeAuth(app *gin.Engine, authcontroller controller.AuthController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/auth")
	{
		routes.POST("/login", authcontroller.Login)
		routes.POST("/register", authcontroller.Register)
		routes.GET("/verify", authcontroller.Verify)
		routes.POST("/send-email", authcontroller.SendVerificationEmail)
		routes.POST("/forget", authcontroller.ForgetPassword)
		routes.POST("/change", authcontroller.ChangePassword)
		routes.GET("/me", middleware.Authenticate(), authcontroller.Me)
		routes.DELETE("/logout", middleware.Authenticate(), authcontroller.Logout)

		// routes.POST("/google", authcontroller.LoginWithGoogle)
	}
}
