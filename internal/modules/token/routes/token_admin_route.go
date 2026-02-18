package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"

	"github.com/gin-gonic/gin"
)

func ServeTokenAdmin(app *gin.Engine, tokenBundleController controller.TokenBundleController, customPricingController controller.CustomPricingController, topupTransactionsController controller.TopupTransactionController, tokenLedgerController controller.TokenLedgerController, summaryController controller.SummaryController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token/admin")
	routes.Use(middleware.Authenticate(), middleware.OnlyAdmin())
	{
		routes.POST("/bundle", tokenBundleController.Create)
		routes.PUT("/bundle/:id", tokenBundleController.Update)
		routes.DELETE("/bundle/:id", tokenBundleController.Delete)

		routes.GET("/custom-pricing", customPricingController.GetCurrent)
		routes.GET("/custom-pricing/history", customPricingController.GetHistory)
		routes.POST("/custom-pricing/toggle-active", customPricingController.ToggleActive)
		routes.POST("/custom-pricing", customPricingController.CreateNewVersion)

		routes.GET("/topup-transactions", topupTransactionsController.GetAll)
		routes.GET("/topup-transactions/:id", topupTransactionsController.GetById)

		routes.GET("/ledger/activity", tokenLedgerController.GetActivity)

		routes.GET("/summary/kpi", summaryController.GetKPI)
		routes.GET("/summary/trend/q", summaryController.GetTrendBySourceType)
		routes.GET("/summary/trend/direction", summaryController.GetTrendDirection)
	}
}
