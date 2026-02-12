package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"

	"github.com/gin-gonic/gin"
)

func ServeToken(
	app *gin.Engine,
	tokenWalletController controller.TokenWalletController,
	tokenBundleController controller.TokenBundleController,
	customPricingController controller.CustomPricingController,
	topupTransactionsController controller.TopupTransactionController,
	tokenLedgerController controller.TokenLedgerController,
	summaryController controller.SummaryController,
	middleware middleware.Middleware) {
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
		routes.POST("/admin/bundle", tokenBundleController.Create)
		routes.PUT("/admin/bundle/:id", tokenBundleController.Update)
		routes.DELETE("/admin/bundle/:id", tokenBundleController.Delete)

		routes.GET("/admin/custom-pricing", customPricingController.GetCurrent)
		routes.GET("/admin/custom-pricing/history", customPricingController.GetHistory)
		routes.POST("/admin/custom-pricing/toggle-active", customPricingController.ToggleActive)
		routes.POST("/admin/custom-pricing", customPricingController.CreateNewVersion)

		routes.GET("/admin/topup-transactions", topupTransactionsController.GetAll)
		routes.GET("/admin/topup-transactions/:id", topupTransactionsController.GetById)

		routes.GET("/admin/ledger/activity", tokenLedgerController.GetActivity)

		routes.GET("/admin/summary/kpi", summaryController.GetKPI)
		routes.GET("/admin/summary/trend/q", summaryController.GetTrendBySourceType)
		routes.GET("/admin/summary/trend/direction", summaryController.GetTrendDirection)
	}
}
