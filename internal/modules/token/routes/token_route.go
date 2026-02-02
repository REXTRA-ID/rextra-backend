package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"

	"github.com/gin-gonic/gin"
)

func ServeToken(app *gin.Engine, tokenWalletController controller.TokenWalletController, tokenBundleController controller.TokenBundleController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token")
	{
		// WALLET
		routes.GET("/wallet/balance", middleware.Authenticate(), tokenWalletController.GetMyBalance)
		routes.GET("/wallet/history", middleware.Authenticate(), tokenWalletController.GetMyHistory)
		routes.GET("/wallet/history/:id", middleware.Authenticate(), tokenWalletController.GetMyHistoryDetail)

		// BUNDLE
		routes.GET("/bundle", tokenBundleController.GetAll)
		routes.GET("/bundle/:id", tokenBundleController.GetByID)
		routes.POST("/bundle", middleware.Authenticate(), middleware.OnlyAdmin(), tokenBundleController.Create)
	}
}
