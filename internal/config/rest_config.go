package config

import (
	"fmt"
	"rextra-backend/db"

	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/api/repository"
	"rextra-backend/internal/api/routes"
	"rextra-backend/internal/api/service"
	"rextra-backend/internal/middleware"
	mailer "rextra-backend/internal/pkg/email"
	myfirebase "rextra-backend/internal/pkg/firebase"

	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type RestConfig struct {
	server *gin.Engine
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
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository                 repository.UserRepository                 = repository.NewUser(db)
		sessionRepository              repository.SessionRepository              = repository.NewSession(db)
		personaRepository              repository.PersonaRepository              = repository.NewPersona(db)
		riasecRepository               repository.RiasecRepository               = repository.NewRiasec(db)
		careerRecommendationRepository repository.CareerRecommendationRepository = repository.NewCareerRecommendation(db)
		assesmentRepository            repository.AssesmentRepository            = repository.NewAssesment(db)
		educationRepository            repository.EducationRepository            = repository.NewEducation(db)

		//=========== (SERVICE) ===========//
		authService                 service.AuthService                 = service.NewAuth(userRepository, sessionRepository, mailerService, firebaseApp.MustGetClient(), db)
		userService                 service.UserService                 = service.NewUser(userRepository, db)
		personaService              service.PersonaService              = service.NewPersona(personaRepository, db)
		riasecService               service.RiasecService               = service.NewRiasec(riasecRepository, db)
		careerRecommendationService service.CareerRecommendationService = service.NewCareerRecommendation(careerRecommendationRepository, db)
		assesmentService            service.AssesmentService            = service.NewAssesment(assesmentRepository, db)
		educationService            service.EducationService            = service.NewEducation(educationRepository, db)

		//=========== (CONTROLLER) ===========//
		authController                 controller.AuthController                 = controller.NewAuth(authService)
		userController                 controller.UserController                 = controller.NewUser(userService)
		personaController              controller.PersonaController              = controller.NewPersona(personaService)
		riasecController               controller.RiasecController               = controller.NewRiasec(riasecService)
		careerRecommendationController controller.CareerRecommendationController = controller.NewCareerRecommendation(careerRecommendationService)
		assesmentController            controller.AssesmentController            = controller.NewAssesment(assesmentService)
		educationController            controller.EducationController            = controller.NewEducation(educationService)
	)

	// Register all routes
	routes.ServeAuth(server, authController, middleware)
	routes.ServeUser(server, userController, middleware)
	routes.ServePersona(server, personaController, middleware)
	routes.ServeRiasec(server, riasecController, middleware)
	routes.ServeCareerRecommendation(server, careerRecommendationController, middleware)
	routes.ServeAssesment(server, assesmentController, middleware)
	routes.ServeEducation(server, educationController, middleware)

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
