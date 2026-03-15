package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/checkout/controller"

	"github.com/gin-gonic/gin"
)

func ServeCheckout(app *gin.Engine, ctrl controller.CheckoutController, mw middleware.Middleware) {
	r := app.Group("/api/v1/checkout")
	r.Use(mw.Authenticate())
	{
		r.GET("/prepare", ctrl.Prepare)
		r.POST("/calculate", ctrl.Calculate)
		r.POST("/initiate", ctrl.Initiate)
		r.GET("/channels", ctrl.GetPaymentChannels)
		r.GET("/transactions", ctrl.GetTransactions)
		r.GET("/transaction/:transactionId", ctrl.GetTransaction)
		r.POST("/repeat/:transactionId", ctrl.Repeat)
		r.POST("/cancel/:transactionId", ctrl.Cancel)
	}
}
