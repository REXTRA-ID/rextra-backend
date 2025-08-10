package config

import (
	"fmt"
	"rextra-backend/db"
	authController "rextra-backend/internal/api/auth/controller"
	authRoutes "rextra-backend/internal/api/auth/routes"
	authService "rextra-backend/internal/api/auth/service"

	userController "rextra-backend/internal/api/user/controller"
	userRepository "rextra-backend/internal/api/user/repository"
	userRoutes "rextra-backend/internal/api/user/routes"
	userService "rextra-backend/internal/api/user/service"
	"rextra-backend/internal/middleware"
	mailer "rextra-backend/internal/pkg/email"
	"rextra-backend/internal/pkg/google/oauth"

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
	middleware := middleware.New(db)

	var (
		//=========== (PACKAGE) ===========//
		mailerService mailer.Mailer = mailer.New()
		oauthService  oauth.Oauth   = oauth.New()
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository userRepository.UserRepository = userRepository.New(db)

		//=========== (SERVICE) ===========//
		authService authService.AuthService = authService.New(userRepository, mailerService, oauthService, db)
		userService userService.UserService = userService.New(userRepository, db)

		//=========== (CONTROLLER) ===========//
		authController authController.AuthController = authController.New(authService)
		userController userController.UserController = userController.New(userService)
	)

	// Register all routes
	authRoutes.Serve(server, authController, middleware)
	userRoutes.Serve(server, userController, middleware)

	return RestConfig{
		server: server,
	}
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")
	if port == "" {
		port = "8998"
	}

	serve := fmt.Sprintf("%s:%s", host, port)
	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
	log.Println("server start on port ", serve)
}
