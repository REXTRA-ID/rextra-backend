package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/payment/controller"

	"github.com/gin-gonic/gin"
)

func ServePaymentTransaction(app *gin.Engine, c controller.PaymentTransactionController, middleware middleware.Middleware) {
	checkout := app.Group("/api/v1/checkout")
	checkout.Use(middleware.Authenticate())
	{
		checkout.GET("/prepare", c.Prepare)
		checkout.POST("/calculate", c.Calculate)
		checkout.POST("/initiate", c.Initiate)
		checkout.POST("/:transactionId/repeat", c.Repeat)
		checkout.POST("/:transactionId/cancel", c.Cancel)
	}

	app.GET("/api/v1/payment-channels", middleware.Authenticate(), c.GetPaymentChannels)

	mytrx := app.Group("/api/v1/my-transactions")
	mytrx.Use(middleware.Authenticate())
	{
		mytrx.GET("", c.GetTransactions)
		mytrx.GET("/:transactionId", c.GetTransaction)
	}
}
