package auth

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth/controller"
	"rextra-backend/internal/modules/auth/repository"
	"rextra-backend/internal/modules/auth/routes"
	"rextra-backend/internal/modules/auth/service"

	mailer "rextra-backend/internal/pkg/email"
	myfirebase "rextra-backend/internal/pkg/firebase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {

	firebaseApp := myfirebase.New()

	var (
		mailerService mailer.Mailer = mailer.New()

		userRepository    repository.UserRepository    = repository.NewUser(db)
		sessionRepository repository.SessionRepository = repository.NewSession(db)

		authService service.AuthService = service.NewAuth(userRepository, sessionRepository, mailerService, firebaseApp.MustGetClient(), db)
		userService service.UserService = service.NewUser(userRepository, db)

		authController controller.AuthController = controller.NewAuth(authService)
		userController controller.UserController = controller.NewUser(userService)
	)

	routes.ServeAuth(server, authController, middleware)
	routes.ServeUser(server, userController, middleware)
}
