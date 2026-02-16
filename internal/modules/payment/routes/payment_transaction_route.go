package routes

import (
	"rextra-backend/internal/modules/payment/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServePaymentTransaction(app *gin.Engine, membershipPlanController controller.PaymentTransactionController, middleware middleware.Middleware) {

	routes := app.Group("/api/v1/transaction")
	{
		routes.POST("/membership", middleware.Authenticate(), membershipPlanController.MakeNewTransactionMembership)
		routes.POST("/token", middleware.Authenticate(), membershipPlanController.MakeNewTransactionTokenStandAlone)
	}
}
