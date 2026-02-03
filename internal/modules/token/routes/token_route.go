package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"

	"github.com/gin-gonic/gin"
)

func ServeToken(app *gin.Engine, tokenWalletController controller.TokenWalletController, tokenBundleController controller.TokenBundleController, customPricingController controller.CustomPricingController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token")
	{
		// WALLET
		routes.GET("/wallet/balance", middleware.Authenticate(), tokenWalletController.GetMyBalance)
		routes.GET("/wallet/history", middleware.Authenticate(), tokenWalletController.GetMyHistory)
		routes.GET("/wallet/history/:id", middleware.Authenticate(), tokenWalletController.GetMyHistoryDetail)

		// BUNDLE
		routes.GET("/bundle", tokenBundleController.GetAll)
		routes.GET("/bundle/:id", tokenBundleController.GetByID)
	}

	routes.Use(middleware.Authenticate(), middleware.OnlyAdmin())
	{
		routes.POST("/bundle", tokenBundleController.Create)
		routes.PUT("/bundle/:id", tokenBundleController.Update)
		routes.DELETE("/bundle/:id", tokenBundleController.Delete)

		routes.GET("/custom-pricing", customPricingController.GetCurrent)
		routes.GET("/custom-pricing/history", customPricingController.GetHistory)
		routes.POST("/custom-pricing/toggle-active", customPricingController.ToggleActive)
		routes.POST("/custom-pricing", customPricingController.CreateNewVersion)
	}
}
