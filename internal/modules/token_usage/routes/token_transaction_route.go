package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token_usage/controller"

	"github.com/gin-gonic/gin"
)

func ServeTokenTransaction(app *gin.Engine, tokenTransactionController controller.TokenTransactionController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token-transactions")
	{
		routes.POST("", middleware.Authenticate(), tokenTransactionController.UseToken)
	}
}
