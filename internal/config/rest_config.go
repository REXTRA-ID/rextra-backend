package config

import (
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/job"
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth"
	"rextra-backend/internal/modules/career_recommendation"
	"rextra-backend/internal/modules/membership"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	membershipService "rextra-backend/internal/modules/membership/service"
	"rextra-backend/internal/modules/payment"
	"rextra-backend/internal/modules/persona"
	"rextra-backend/internal/modules/poin"
	"rextra-backend/internal/modules/riasec"
	"rextra-backend/internal/modules/token"
	mailer "rextra-backend/internal/pkg/email"
	"rextra-backend/internal/pkg/tripay"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	// Initialize all modules
	auth.InitModule(server, db, middleware, mailerService)
	membership.InitModule(server, db, middleware)
	payment.InitModule(server, db, middleware, &tripayClient)
	token.InitModule(server, db, middleware)
	poin.InitModule(server, db, middleware)
	persona.InitModule(server, db, middleware)
	riasec.InitModule(server, db, middleware)
	career_recommendation.InitModule(server, db, middleware)

	// Cronjobs - need to create services for jobs
	c := cron.New(cron.WithLogger(cron.DefaultLogger))

	// Create repositories and services for cronjobs
	membershipRepository := membershipRepo.NewMembershipRepository(db)
	membershipPlanRepository := membershipRepo.NewMembershipPlanRepository(db)
	membershipSvc := membershipService.NewMembershipService(membershipRepository, membershipPlanRepository)

	// Note: Cronjobs are temporarily disabled until we refactor job package to use module services
	// tokenTransactionRepository := tokenRepo.NewTokenTransactionRepository(db)
	// tokenUsageHistoryRepository := tokenRepo.NewTokenUsageHistoryRepository(db)
	// tokenTransactionSvc := tokenService.NewTokenTransactionService(membershipRepository, tokenTransactionRepository, tokenUsageHistoryRepository, db)

	expireMembershipjob := &job.ExpireMembershipJob{
		MembershipService: membershipSvc,
	}

	// refillTokenJob := &job.RefillTokenJob{
	// 	TokenTransactionService: tokenTransactionSvc,
	// }

	// c.AddJob("0 0 1 * *", refillTokenJob)
	c.AddJob("0 0 * * *", expireMembershipjob)

	c.Start()

	return RestConfig{
		server: server,
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
