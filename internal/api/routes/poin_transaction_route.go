package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServePoinTransaction(app *gin.Engine, poinTransactionController controller.PoinTransactionController, middlewre middleware.Middleware) {
	routes := app.Group("/api/v1/poin-transaction")
	{
		routes.POST("")
		routes.GET("")
	}
}
