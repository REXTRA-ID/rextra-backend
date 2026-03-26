package config

import (
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/job"
	"rextra-backend/internal/middleware"
	actioncategory "rextra-backend/internal/modules/action_category"
	"rextra-backend/internal/modules/auth"
	"rextra-backend/internal/modules/entitlement"
	"rextra-backend/internal/modules/feature"
	"rextra-backend/internal/modules/membership"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	membershipService "rextra-backend/internal/modules/membership/service"
	"rextra-backend/internal/modules/payment"
	"rextra-backend/internal/modules/persona"
	"rextra-backend/internal/modules/poin"
	tokenRepo "rextra-backend/internal/modules/token_usage/repository"
	tokenService "rextra-backend/internal/modules/token_usage/service"
	"rextra-backend/internal/pkg/tripay"

	"log"
	"os"
	"rextra-backend/internal/modules/access_check"
	kenalidiri "rextra-backend/internal/modules/kenali_diri"
	token "rextra-backend/internal/modules/token"
	tokenModuleRepo "rextra-backend/internal/modules/token/repository"

	"rextra-backend/internal/pkg/cache"
	"rextra-backend/internal/pkg/export"
	myfirebase "rextra-backend/internal/pkg/firebase"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type RestConfig struct {
	server       *gin.Engine
	cacheService cache.CacheService
}

func NewRest() RestConfig {
	db := db.New()

	// Mode
	mode := os.Getenv("APP_MODE")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := gin.Default()
	server := NewRouter(app)
	tripayClient := tripay.NewTripayClient()

	firebaseApp := myfirebase.New()
	middleware := middleware.New(db, firebaseApp.MustGetClient())

	tokenWalletRepo := tokenModuleRepo.NewTokenWalletRepository(db)
	tokenLedgerRepo := tokenModuleRepo.NewTokenLedgerRepository(db)
	accessCheckSvc := access_check.NewAccessCheckService(db, tokenWalletRepo, tokenLedgerRepo)
	middleware.SetAccessCheckService(accessCheckSvc)

	// Cronjobs
	c := cron.New(cron.WithLogger(cron.DefaultLogger))

	membershipRepository := membershipRepo.NewMembershipRepository(db)
	membershipPlanRepository := membershipRepo.NewMembershipPlanRepository(db)
	membershipSvc := membershipService.NewMembershipService(membershipRepository, membershipPlanRepository)

	tokenTransactionRepository := tokenRepo.NewTokenTransactionRepository(db)
	tokenUsageHistoryRepository := tokenRepo.NewTokenUsageHistoryRepository(db)
	tokenTransactionSvc := tokenService.NewTokenTransactionService(membershipRepository, tokenTransactionRepository, tokenUsageHistoryRepository, db)

	expireMembershipjob := &job.ExpireMembershipJob{
		MembershipService: membershipSvc,
	}

	refillTokenJob := &job.RefillTokenJob{
		TokenTransactionService: tokenTransactionSvc,
	}

	c.AddJob("0 0 * * *", expireMembershipjob)
	c.AddJob("0 0 * * *", refillTokenJob)
	c.Start()

	cacheService := cache.New()
	var exportService export.ExportService = export.New()

	// Initialize all modules
	auth.InitModule(server, db, middleware)
	persona.InitModule(server, db, middleware)
	kenalidiri.InitModule(server, db, middleware, cacheService, exportService)
	token.InitModule(server, db, middleware)
	poin.InitModule(server, db, middleware)
	payment.InitModule(server, db, middleware, &tripayClient)
	membership.InitModule(server, db, middleware)

	// Hak akses modules
	actioncategory.InitModule(server, db, middleware)
	feature.InitModule(server, db, middleware)
	entitlement.InitModule(server, db, middleware)

	return RestConfig{
		server:       server,
		cacheService: cacheService,
	}
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")
	if port == "" {
		port = "8000"
	}

	serve := fmt.Sprintf("%s:%s", host, port)
	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
	log.Println("server start on port ", serve)
}

func (ap *RestConfig) Close() error {
	if ap.cacheService != nil {
		if closer, ok := ap.cacheService.(interface{ Close() error }); ok {
			return closer.Close()
		}
	}
	return nil
}
