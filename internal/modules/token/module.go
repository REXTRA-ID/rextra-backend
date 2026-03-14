package token

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/token/routes"
	"rextra-backend/internal/modules/token/service"

	membershipRepo "rextra-backend/internal/modules/membership/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repository
	membershipRepository := membershipRepo.NewUserMembershipRepository(db)
	tokenTransactionRepository := repository.NewTokenTransactionRepository(db)
	tokenUsageHistoryRepository := repository.NewTokenUsageHistoryRepository(db)

	// Service
	tokenTransactionService := service.NewTokenTransactionService(
		membershipRepository,
		tokenTransactionRepository,
		tokenUsageHistoryRepository,
		db,
	)

	// Controller
	tokenTransactionController := controller.NewTokenTransactionController(tokenTransactionService)

	// Routes
	routes.ServeTokenTransaction(server, tokenTransactionController, middleware)
}
