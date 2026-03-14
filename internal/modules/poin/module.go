package poin

import (
	"rextra-backend/internal/middleware"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/poin/controller"
	"rextra-backend/internal/modules/poin/repository"
	"rextra-backend/internal/modules/poin/routes"
	"rextra-backend/internal/modules/poin/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repositories
	poinTransactionRepository := repository.NewPoinTransactionsRepository(db)

	// Cross-module repository
	membershipRepository := membershipRepo.NewMembershipRepository(db)

	// Service
	poinTransactionService := service.NewPoinTransactionService(
		membershipRepository,
		poinTransactionRepository,
		db,
	)

	// Controller
	poinTransactionController := controller.NewPoinTransactionController(poinTransactionService)

	// Routes
	routes.ServePoinTransaction(server, poinTransactionController, middleware)
}
