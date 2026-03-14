package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/payment/controller"

	"github.com/gin-gonic/gin"
)

func ServePaymentTransaction(
	app *gin.Engine,
	ctrl controller.PaymentTransactionController,
	mw middleware.Middleware,
) {
	public := app.Group("/api/v1/payment")
	{
		public.GET("/channels", ctrl.GetPaymentChannels)
		public.POST("/callback/tripay", ctrl.HandleTripayCallback)
		public.POST("/simulate/:transactionId", ctrl.SimulateTripayPayment)
	}

	userPayment := app.Group("/api/v1/payment")
	userPayment.Use(mw.Authenticate())
	{
		userPayment.POST("/transaction", ctrl.CreateTransaction)
		userPayment.POST("/calculate", ctrl.CalculatePrice)
	}

	myRoutes := app.Group("/api/v1/my")
	myRoutes.Use(mw.Authenticate())
	{
		myRoutes.GET("/payment", ctrl.GetMyTransactions)
		myRoutes.GET("/payment/:transactionId", ctrl.GetMyTransactionDetail)
	}

	adminRoutes := app.Group("/api/v1/admin/payment")
	adminRoutes.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		adminRoutes.GET("", ctrl.GetAllTransactions)
		adminRoutes.GET("/:transactionId", ctrl.GetTransactionDetail)
		adminRoutes.PUT("/:transactionId/cancel", ctrl.CancelTransaction)
	}
}
