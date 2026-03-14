package payment

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/payment/controller"
	payrepo "rextra-backend/internal/modules/payment/repository"
	"rextra-backend/internal/modules/payment/routes"
	"rextra-backend/internal/modules/payment/service"
	"rextra-backend/internal/pkg/tripay"

	discountRepo "rextra-backend/internal/modules/promo/repository"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	poinRepo "rextra-backend/internal/modules/poin/repository"
	subCycleRepo "rextra-backend/internal/modules/subscription/repository"
	tokenRepo "rextra-backend/internal/modules/token/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware, tripayClient *tripay.TripayClient) {
	paymentRepository := payrepo.NewPaymentTransactionRepository(db)

	membershipRepository := membershipRepo.NewUserMembershipRepository(db)
	planRepository := membershipRepo.NewMembershipPlanRepository(db)
	durationRepository := membershipRepo.NewPlanDurationRepository(db)
	cycleRepository := subCycleRepo.NewSubscriptionCycleRepository(db)
	discountRepository := discountRepo.NewDiscountRepository(db)
	redemptionRepository := discountRepo.NewDiscountRedemptionRepository(db)
	tokenLedgerRepository := tokenRepo.NewTokenLedgerRepository(db)
	tokenWalletRepository := tokenRepo.NewTokenWalletRepository(db)
	topupTrxRepository := tokenRepo.NewTopupTransactionRepository(db)
	poinRepository := poinRepo.NewPoinTransactionsRepository(db)
	accessMappingRepository := membershipRepo.NewDurationAccessMappingRepository(db)
	quotaRepository := membershipRepo.NewUserEntitlementQuotaRepository(db)

	paymentService := service.NewPaymentTransactionService(
		paymentRepository,
		membershipRepository,
		planRepository,
		durationRepository,
		cycleRepository,
		discountRepository,
		redemptionRepository,
		redemptionRepository,
		tokenLedgerRepository,
		topupTrxRepository,
		tokenWalletRepository,
		poinRepository,
		accessMappingRepository,
		quotaRepository,
		*tripayClient,
		db,
	)

	paymentController := controller.NewPaymentTransactionController(paymentService)

	routes.ServePaymentTransaction(server, paymentController, mw)
}
