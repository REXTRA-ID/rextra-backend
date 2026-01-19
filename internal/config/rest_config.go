package config

import (
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/api/controller"
	kdcontroller "rextra-backend/internal/api/kenali_diri/controller"
	kdrepo "rextra-backend/internal/api/kenali_diri/repository"
	kdroutes "rextra-backend/internal/api/kenali_diri/routes"
	kdservice "rextra-backend/internal/api/kenali_diri/service"
	"rextra-backend/internal/api/repository"
	"rextra-backend/internal/api/routes"
	"rextra-backend/internal/api/service"
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/pkg/cache"
	mailer "rextra-backend/internal/pkg/email"
	"rextra-backend/internal/pkg/export"
	myfirebase "rextra-backend/internal/pkg/firebase"

	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// initCache initializes cache service with Redis if available, fallback to memory cache
func initCache() cache.CacheService {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0 // Default to DB 0

	// Try to initialize Redis cache
	if redisHost != "" {
		redisCache, err := cache.NewRedis(redisHost, redisPort, redisPassword, redisDB)
		if err != nil {
			log.Printf("⚠️  Failed to connect to Redis at %s:%s, falling back to memory cache: %v", redisHost, redisPort, err)
			return cache.NewMemory()
		}
		log.Printf("✅ Redis cache initialized successfully at %s:%s", redisHost, redisPort)
		return redisCache
	}

	// If REDIS_HOST not set, use memory cache
	log.Println("ℹ️  REDIS_HOST not configured, using memory cache")
	return cache.NewMemory()
}

type RestConfig struct {
	server       *gin.Engine
	cacheService cache.CacheService
}

func NewRest() RestConfig {
	db := db.New()
	app := gin.Default()
	server := NewRouter(app)
	firebaseApp := myfirebase.New()
	middleware := middleware.New(db, firebaseApp.MustGetClient())

	var (
		//=========== (PACKAGE) ===========//
		mailerService mailer.Mailer = mailer.New()
		exportService export.ExportService = export.New()
		cacheService  cache.CacheService  = initCache()
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository                 repository.UserRepository                 = repository.NewUser(db)
		sessionRepository              repository.SessionRepository              = repository.NewSession(db)
		personaRepository              repository.PersonaRepository              = repository.NewPersona(db)
		riasecRepository               repository.RiasecRepository               = repository.NewRiasec(db)
		careerRecommendationRepository repository.CareerRecommendationRepository = repository.NewCareerRecommendation(db)
		assesmentRepository            repository.AssesmentRepository            = repository.NewAssesment(db)
		kenalidiriHistoryRepository    kdrepo.KenalidiriHistoryRepository        = kdrepo.NewKenalidiriHistory(db)
		kenalidiriCategoryRepository   kdrepo.KenalidiriCategoryRepository       = kdrepo.NewKenalidiriCategory(db, cacheService)
		kenalidiriRiasecCodeRepository kdrepo.RiasecCodeRepository               = kdrepo.NewRiasecCode(db, cacheService)
		testSessionRepository          kdrepo.TestSessionRepository              = kdrepo.NewTestSession(db)
		kenalidiriRiasecRepository     kdrepo.RiasecRepository                   = kdrepo.NewRiasec(db)
		ikigaiRepository               kdrepo.IkigaiRepository                   = kdrepo.NewIkigai(db)
		recommendationRepository       kdrepo.RecommendationRepository           = kdrepo.NewRecommendation(db)
		feedbackRepository             kdrepo.FeedbackRepository                 = kdrepo.NewFeedback(db)

		//=========== (SERVICE) ===========//
		authService                 service.AuthService                 = service.NewAuth(userRepository, sessionRepository, mailerService, firebaseApp.MustGetClient(), db)
		userService                 service.UserService                 = service.NewUser(userRepository, db)
		personaService              service.PersonaService              = service.NewPersona(personaRepository, db)
		riasecService               service.RiasecService               = service.NewRiasec(riasecRepository, db)
		careerRecommendationService service.CareerRecommendationService = service.NewCareerRecommendation(careerRecommendationRepository, db)
		assesmentService            service.AssesmentService            = service.NewAssesment(assesmentRepository, db)
		kenalidiriAdminService      kdservice.KenalidiriAdminService    = kdservice.NewKenalidiriAdmin(
			kenalidiriHistoryRepository,
			kenalidiriCategoryRepository,
			kenalidiriRiasecCodeRepository,
			testSessionRepository,
			kenalidiriRiasecRepository,
			ikigaiRepository,
			recommendationRepository,
			feedbackRepository,
			exportService,
			cacheService,
			db,
		)

		//=========== (CONTROLLER) ===========//
		authController                 controller.AuthController                 = controller.NewAuth(authService)
		userController                 controller.UserController                 = controller.NewUser(userService)
		personaController              controller.PersonaController              = controller.NewPersona(personaService)
		riasecController               controller.RiasecController               = controller.NewRiasec(riasecService)
		careerRecommendationController controller.CareerRecommendationController = controller.NewCareerRecommendation(careerRecommendationService)
		assesmentController            controller.AssesmentController            = controller.NewAssesment(assesmentService)
		kenalidiriAdminController      kdcontroller.KenalidiriAdminController    = kdcontroller.NewKenalidiriAdmin(kenalidiriAdminService)
	)

	// Register all routes
	routes.ServeAuth(server, authController, middleware)
	routes.ServeUser(server, userController, middleware)
	routes.ServePersona(server, personaController, middleware)
	routes.ServeRiasec(server, riasecController, middleware)
	routes.ServeCareerRecommendation(server, careerRecommendationController, middleware)
	routes.ServeAssesment(server, assesmentController, middleware)
	kdroutes.ServeKenalidiriAdmin(server, kenalidiriAdminController, middleware)

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
