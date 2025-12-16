package config

import (
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/api/repository"
	"rextra-backend/internal/api/routes"
	"rextra-backend/internal/api/service"
	"rextra-backend/internal/job"
	"rextra-backend/internal/middleware"
	mailer "rextra-backend/internal/pkg/email"
	"rextra-backend/payment_handler/midtrans"

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
	/* untuk sementara */
	// firebaseApp := myfirebase.New()
	middleware := middleware.New(db)
	// xenditService                 := xnd.NewXenditService() xendit service
	midtransService := midtrans.NewMidtransClient()
	var (
		//=========== (PACKAGE) ===========//
		mailerService mailer.Mailer = mailer.New()
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository                 repository.UserRepository                 = repository.NewUser(db)
		sessionRepository              repository.SessionRepository              = repository.NewSession(db)
		personaRepository              repository.PersonaRepository              = repository.NewPersona(db)
		riasecRepository               repository.RiasecRepository               = repository.NewRiasec(db)
		careerRecommendationRepository repository.CareerRecommendationRepository = repository.NewCareerRecommendation(db)
		membershipRepository           repository.MembershipRepository           = repository.NewMembershipRepository(db)
		membershipPlanRepository       repository.MembershipPlanRepository       = repository.NewMembershipPlanRepository(db)
		membershipDurationRepository   repository.MembershipDurationRepository   = repository.NewMembershipDurationRepository(db)
		paymentTransactionRepository   repository.PaymentTransactionsRepository  = repository.NewPaymentTransactionRepository(db)
		tokenUsageHistoryRepository    repository.TokenUsageHistoryRepository    = repository.NewTokenUsageHistoryRepository(db)
		tokenTransactionRepository     repository.TokenTransactionRepository     = repository.NewTokenTransactionRepository(db)
		poinTransactionRepository      repository.PoinTransactionsRepository     = repository.NewPoinTransactionsRepository(db)
		promoCodeRepository            repository.PromoCodeRepository            = repository.NewPromoCodesRepository(db)
		promoCodeUsageRepository       repository.PromoCodeUsageRepository       = repository.NewPromoCodeUsageRepository(db)

		//=========== (SERVICE) ===========//
		authService                 service.AuthService                 = service.NewAuth(userRepository, membershipRepository, membershipDurationRepository, membershipPlanRepository, sessionRepository, mailerService, nil, db)
		userService                 service.UserService                 = service.NewUser(userRepository, db)
		personaService              service.PersonaService              = service.NewPersona(personaRepository, db)
		riasecService               service.RiasecService               = service.NewRiasec(riasecRepository, db)
		careerRecommendationService service.CareerRecommendationService = service.NewCareerRecommendation(careerRecommendationRepository, db)
		membershipService           service.MembershipService           = service.NewMembershipService(membershipRepository, membershipPlanRepository)
		membershipPlanService       service.MembershipPlanService       = service.NewMembershipPlanService(membershipPlanRepository, db)
		membershipDurationService   service.MembershipDurationService   = service.NewMembershipDurationService(membershipDurationRepository, db)
		tokenTransactionService     service.TokenTransactionService     = service.NewTokenTransactionService(membershipRepository, tokenTransactionRepository, tokenUsageHistoryRepository, db)
		poinTransactionService      service.PoinTransactionService      = service.NewPoinTransactionService(membershipRepository, poinTransactionRepository, db)
		paymentTransactionService   service.PaymentTransactionService   = service.NewPaymentTransactionService(tokenTransactionRepository, poinTransactionRepository, midtransService, paymentTransactionRepository, promoCodeRepository, promoCodeUsageRepository, membershipPlanRepository, membershipDurationRepository, membershipRepository, db)

		//=========== (CONTROLLER) ===========//
		authController                 controller.AuthController                 = controller.NewAuth(authService)
		userController                 controller.UserController                 = controller.NewUser(userService)
		personaController              controller.PersonaController              = controller.NewPersona(personaService)
		riasecController               controller.RiasecController               = controller.NewRiasec(riasecService)
		careerRecommendationController controller.CareerRecommendationController = controller.NewCareerRecommendation(careerRecommendationService)
		membershipPlanController       controller.MembershipPlanController       = controller.NewMembership(membershipPlanService)
		membershipDurationController   controller.MembershipDurationController   = controller.NewMembershipDurationController(membershipDurationService)
		paymentTransactionController   controller.PaymentTransactionController   = controller.NewPaymentTransactionController(paymentTransactionService)
		tokenTransactionController     controller.TokenTransactionController     = controller.NewTokenTransactionController(tokenTransactionService)
		poinTransactionController      controller.PoinTransactionController      = controller.NewPoinTransactionController(poinTransactionService)
	)

	// Cronjobs
	c := cron.New(cron.WithLogger(cron.DefaultLogger))

	expireMembershipjob := &job.ExpireMembershipJob{
		MembershipService: membershipService,
	}

	refillTokenJob := &job.RefillTokenJob{
		TokenTransactionService: tokenTransactionService,
	}

	c.AddJob("0 0 1 * *", refillTokenJob)
	c.AddJob("0 0 * * *", expireMembershipjob)

	c.Start()

	// Register all routes
	routes.ServeAuth(server, authController, middleware)
	routes.ServeUser(server, userController, middleware)
	routes.ServePersona(server, personaController, middleware)
	routes.ServeRiasec(server, riasecController, middleware)
	routes.ServeCareerRecommendation(server, careerRecommendationController, middleware)
	routes.ServeMembershipPlan(server, membershipPlanController, middleware)
	routes.ServeMemberDuration(server, membershipDurationController, middleware)
	routes.ServePaymentTransaction(server, paymentTransactionController, middleware)
	routes.ServeTokenTransaction(server, tokenTransactionController, middleware)
	routes.ServePoinTransaction(server, poinTransactionController, middleware)

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
