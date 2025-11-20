package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServePaymentTransaction(app *gin.Engine, membershipPlanController controller.PaymentTransactionController, middleware middleware.Middleware) {

	routes := app.Group("/api/v1/transaction")
	{
		routes.POST("/membership")
	}
}
