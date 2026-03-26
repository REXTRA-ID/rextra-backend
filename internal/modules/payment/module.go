package payment

import (
	"rextra-backend/internal/middleware"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/payment/controller"
	"rextra-backend/internal/modules/payment/repository"
	"rextra-backend/internal/modules/payment/routes"
	"rextra-backend/internal/modules/payment/service"
	userRepo "rextra-backend/internal/modules/user/repository"
	"rextra-backend/internal/pkg/tripay"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware, tripayClient *tripay.TripayClient) {
	// Repositories
	paymentTransactionRepository := repository.NewPaymentTransactionRepository(db)

	// Cross-module repositories
	membershipPlanRepository := membershipRepo.NewMembershipPlanRepository(db)
	membershipDurationRepository := membershipRepo.NewMembershipDurationRepository(db)
	userRepository := userRepo.NewUser(db)
	promoFacade := service.NewPromoFacade()

	// Service
	paymentTransactionService := service.NewPaymentTransactionService(
		paymentTransactionRepository,
		membershipPlanRepository,
		membershipDurationRepository,
		promoFacade,
		userRepository,
		tripayClient,
		db,
	)

	// Controller
	paymentTransactionController := controller.NewPaymentTransactionController(paymentTransactionService)

	// Routes
	routes.ServePaymentTransaction(server, paymentTransactionController, middleware)
}
