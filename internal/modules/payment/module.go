package payment

import (
	"rextra-backend/internal/middleware"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/payment/controller"
	"rextra-backend/internal/modules/payment/repository"
	"rextra-backend/internal/modules/payment/routes"
	"rextra-backend/internal/modules/payment/service"
	poinRepo "rextra-backend/internal/modules/poin/repository"
	promoCodeRepo "rextra-backend/internal/modules/promo_code/repository"
	tokenRepo "rextra-backend/internal/modules/token_usage/repository"
	"rextra-backend/internal/pkg/tripay"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware, tripay *tripay.TripayClient) {
	// Repositories
	paymentTransactionRepository := repository.NewPaymentTransactionRepository(db)

	// Cross-module repositories
	tokenTransactionRepository := tokenRepo.NewTokenTransactionRepository(db)
	poinTransactionRepository := poinRepo.NewPoinTransactionsRepository(db)
	promoCodeRepository := promoCodeRepo.NewPromoCodesRepository(db)
	promoCodeUsageRepository := promoCodeRepo.NewPromoCodeUsageRepository(db)
	membershipPlanRepository := membershipRepo.NewMembershipPlanRepository(db)
	membershipDurationRepository := membershipRepo.NewMembershipDurationRepository(db)
	membershipRepository := membershipRepo.NewMembershipRepository(db)

	// Service
	paymentTransactionService := service.NewPaymentTransactionService(
		tokenTransactionRepository,
		poinTransactionRepository,
		*tripay,
		paymentTransactionRepository,
		promoCodeRepository,
		promoCodeUsageRepository,
		membershipPlanRepository,
		membershipDurationRepository,
		membershipRepository,
		db,
	)

	// Controller
	paymentTransactionController := controller.NewPaymentTransactionController(paymentTransactionService)

	// Routes
	routes.ServePaymentTransaction(server, paymentTransactionController, middleware)
}
