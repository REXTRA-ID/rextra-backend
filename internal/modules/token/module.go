package token

import (
	"rextra-backend/internal/middleware"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/token/controller"
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/token/routes"
	"rextra-backend/internal/modules/token/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repositories
	tokenTransactionRepository := repository.NewTokenTransactionRepository(db)
	tokenUsageHistoryRepository := repository.NewTokenUsageHistoryRepository(db)

	// Cross-module repository
	membershipRepository := membershipRepo.NewMembershipRepository(db)

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
