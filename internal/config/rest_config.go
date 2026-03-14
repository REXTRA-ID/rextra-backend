package config

import (
	"context"
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/access_check"
	action_category "rextra-backend/internal/modules/action_category"
	"rextra-backend/internal/modules/auth"
	"rextra-backend/internal/modules/checkout"
	"rextra-backend/internal/modules/entitlement"
	"rextra-backend/internal/modules/feature"
	"rextra-backend/internal/modules/membership"
	membershipRepository "rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/payment"
	"rextra-backend/internal/modules/pengaturan"
	"rextra-backend/internal/modules/poin"
	"rextra-backend/internal/modules/promo"
	promoRepo "rextra-backend/internal/modules/promo/repository"
	promoService "rextra-backend/internal/modules/promo/service"
	"rextra-backend/internal/modules/subscription"
	"rextra-backend/internal/modules/token"
	tokenRepo "rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/transaction_history"
	"rextra-backend/internal/modules/user"
	"rextra-backend/internal/modules/user_entitlement"
	mailer "rextra-backend/internal/pkg/email"
	"rextra-backend/internal/pkg/tripay"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestConfig struct {
	server *gin.Engine
}

func NewRest() RestConfig {
	db := db.New()
	app := gin.Default()
	server := NewRouter(app)
	middleware := middleware.New(db)
	tripayClient := tripay.NewTripayClient()

	var (
		mailerService mailer.Mailer = mailer.New()
	)

	// 1. Initialize Base & Master Data Modules
	auth.InitModule(server, db, middleware, mailerService)
	user.InitModule(server, db, middleware)
	action_category.InitModule(server, db, middleware)
	feature.InitModule(server, db, middleware)
	entitlement.InitModule(server, db, middleware)
	
	// 2. Initialize Shared Logic Modules
	token.InitModule(server, db, middleware)
	poin.InitModule(server, db, middleware)
	
	// 3. Initialize Settings & Promo
	pengaturan.InitModule(server, db, middleware)
	promo.InitModule(server, db, middleware)

	// 4. Initialize Core Business Modules
	membership.InitModule(server, db, middleware)
	payment.InitModule(server, db, middleware, &tripayClient)
	subscription.InitModule(server, db, middleware)
	
	// 5. Initialize User Services
	access_check.NewAccessCheckService(db, tokenRepo.NewTokenWalletRepository(db), tokenRepo.NewTokenLedgerRepository(db))
	
	discountRepository := promoRepo.NewDiscountRepository(db)
	redemptionRepository := promoRepo.NewDiscountRedemptionRepository(db)
	
	// Adapter for PromoService to read plan duration
	planDurationRepo := membershipRepository.NewPlanDurationRepository(db)
	planDurationPromoAdapter := &planDurationPromoAdapter{repo: planDurationRepo}
	promoSvc := promoService.NewDiscountService(discountRepository, redemptionRepository, planDurationPromoAdapter, db)

	checkout.InitModule(server, db, middleware, &tripayClient, promoSvc)
	transaction_history.InitModule(server, db, middleware)
	user_entitlement.InitModule(server, db, middleware)

	return RestConfig{
		server: server,
	}
}

type planDurationPromoAdapter struct {
	repo membershipRepository.PlanDurationRepository
}

func (a *planDurationPromoAdapter) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (promoService.PlanDurationInfo, error) {
	duration, err := a.repo.GetByID(ctx, tx, id)
	if err != nil { return promoService.PlanDurationInfo{}, err }
	planName := ""
	if duration.Plan != nil { planName = string(duration.Plan.PlanName) }
	return promoService.PlanDurationInfo{
		ID:             duration.ID,
		PlanID:         duration.PlanID,
		PlanName:       planName,
		FinalPrice:     duration.FinalPrice,
		DurationMonths: duration.DurationMonths,
	}, nil
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	serve := fmt.Sprintf(":%s", port)
	if host := os.Getenv("APP_HOST"); host != "" {
		serve = fmt.Sprintf("%s:%s", host, port)
	}

	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
}
