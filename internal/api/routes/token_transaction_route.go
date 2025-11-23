package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeTokenTransaction(app *gin.Engine, tokenTransactionController controller.TokenTransactionController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/token-transactions")
	{
		routes.POST("")
	}
}
