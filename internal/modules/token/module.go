package persona

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/token/controller"
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/token/routes"
	"rextra-backend/internal/modules/token/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	var (
		tokenBundleRepository   repository.TokenBundlePackageRepository = repository.NewTokenBundlePackageRepository(db)
		tokenLedgerRepository   repository.TokenLedgerRepository        = repository.NewTokenLedgerRepository(db)
		tokenWalletRepository   repository.TokenWalletRepository        = repository.NewTokenWalletRepository(db)
		customPricingRepository repository.CustomPricingRepository      = repository.NewCustomPricingRepository(db)
		topUpRepository         repository.TopupTransactionRepository   = repository.NewTopupTransactionRepository(db)

		walletService        service.WalletService           = service.NewWalletService(tokenWalletRepository, tokenLedgerRepository, db)
		bundleService        service.BundleService           = service.NewBundleService(tokenBundleRepository, db)
		customPricingService service.CustomPricingService    = service.NewCustomPricingService(customPricingRepository, db)
		topUpService         service.TopupTransactionService = service.NewTopupTransactionService(topUpRepository, db)
		tokenLedgerService   service.TokenLedgerService      = service.NewTokenLedgerService(tokenLedgerRepository, db)
		summaryService       service.TokenSummaryService     = service.NewTokenSummaryService(tokenLedgerRepository, topUpRepository, db)
		paymentService       service.PaymentService          = service.NewPaymentService(topUpRepository, tokenLedgerRepository)

		tokenBundleController   controller.TokenBundleController      = controller.NewTokenBundleController(bundleService)
		tokenWalletController   controller.TokenWalletController      = controller.NewTokenWalletController(walletService)
		customPricingController controller.CustomPricingController    = controller.NewCustomPricingController(customPricingService)
		topUpController         controller.TopupTransactionController = controller.NewTopupTransactionController(topUpService)
		tokenLedgerController   controller.TokenLedgerController      = controller.NewTokenLedgerController(tokenLedgerService)
		tokenSummaryController  controller.SummaryController          = controller.NewSummaryController(summaryService)
		paymentController       controller.PaymentController          = controller.NewPaymentController(paymentService)
	)

	routes.ServeToken(server, tokenWalletController, tokenBundleController, paymentController, middleware)
	routes.ServeTokenAdmin(server, tokenBundleController, customPricingController, topUpController, tokenLedgerController, tokenSummaryController, middleware)
}
