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
		tokenBundleRepository repository.TokenBundlePackageRepository = repository.NewTokenBundlePackageRepository(db)
		tokenLedgerRepository repository.TokenLedgerRepository        = repository.NewTokenLedgerRepository(db)
		tokenWalletRepository repository.TokenWalletRepository        = repository.NewTokenWalletRepository(db)
		// customPricingRepository repository.CustomPricingRepository = repository.NewCustomPricingRepository(db)
		// topUpRepository repository.TopupTransactionRepository = repository.NewTopupTransactionRepository(db)

		walletService service.WalletService = service.NewWalletService(tokenWalletRepository, tokenLedgerRepository, db)
		bundleService service.BundleService = service.NewBundleService(tokenBundleRepository, db)

		tokenBundleController controller.TokenBundleController = controller.NewTokenBundleController(bundleService)
		tokenWalletController controller.TokenWalletController = controller.NewTokenWalletController(walletService)
	)

	routes.ServeToken(server, tokenWalletController, tokenBundleController, middleware)
}
