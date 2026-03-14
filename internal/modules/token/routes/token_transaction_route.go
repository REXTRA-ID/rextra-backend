package routes

import (
	"rextra-backend/internal/modules/token/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeTokenTransaction(app *gin.Engine, tokenTransactionController controller.TokenTransactionController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token-transactions")
	{
		routes.POST("", middleware.Authenticate(), tokenTransactionController.UseToken)
	}
}
