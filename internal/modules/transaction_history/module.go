package transaction_history

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/transaction_history/controller"
	"rextra-backend/internal/modules/transaction_history/repository"
	"rextra-backend/internal/modules/transaction_history/routes"
	"rextra-backend/internal/modules/transaction_history/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	repo := repository.NewHistoryRepository(db)
	svc := service.NewHistoryService(repo)
	ctrl := controller.NewHistoryController(svc)
	routes.ServeHistory(server, ctrl, middleware)
}
