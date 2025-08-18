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
		userRepository    repository.UserRepository    = repository.NewUser(db)
		sessionRepository repository.SessionRepository = repository.NewSession(db)
		personaRepository repository.PersonaRepository = repository.NewPersona(db)
		riasecRepository  repository.RiasecRepository  = repository.NewRiasec(db)

		//=========== (SERVICE) ===========//
		authService    service.AuthService    = service.NewAuth(userRepository, sessionRepository, mailerService, firebaseApp.MustGetClient(), db)
		userService    service.UserService    = service.NewUser(userRepository, db)
		personaService service.PersonaService = service.NewPersona(personaRepository, db)
		riasecService  service.RiasecService  = service.NewRiasec(riasecRepository, db)

		//=========== (CONTROLLER) ===========//
		authController    controller.AuthController    = controller.NewAuth(authService)
		userController    controller.UserController    = controller.NewUser(userService)
		personaController controller.PersonaController = controller.NewPersona(personaService)
		riasecController  controller.RiasecController  = controller.NewRiasec(riasecService)
	)

	// Register all routes
	routes.ServeAuth(server, authController, middleware)
	routes.ServeUser(server, userController, middleware)
	routes.ServePersona(server, personaController, middleware)
	routes.ServeRiasec(server, riasecController, middleware)

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
