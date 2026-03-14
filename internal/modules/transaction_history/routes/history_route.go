package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/transaction_history/controller"

	"github.com/gin-gonic/gin"
)

func ServeHistory(server *gin.Engine, ctrl controller.HistoryController, m middleware.Middleware) {
	r := server.Group("/api/v1/transaction-history", m.Authenticate())
	{
		r.GET("/transactions", ctrl.GetTransactions)
	}
}
