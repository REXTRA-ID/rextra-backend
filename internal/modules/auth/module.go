package auth

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth/controller"
	auth "rextra-backend/internal/modules/auth/repository"
	"rextra-backend/internal/modules/auth/routes"
	"rextra-backend/internal/modules/auth/service"
	membership "rextra-backend/internal/modules/membership/repository"
	user "rextra-backend/internal/modules/user/repository"
	mailer "rextra-backend/internal/pkg/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware, mailerService mailer.Mailer) {
	// Repository
	userRepository := user.NewUserRepository(db)
	membershipRepository := membership.NewUserMembershipRepository(db)
	membershipPlanRepository := membership.NewMembershipPlanRepository(db)
	sessionRepository := auth.NewSession(db)

	// Service
	authService := service.NewAuth(
		userRepository,
		membershipRepository,
		membershipPlanRepository,
		sessionRepository,
		mailerService,
		nil, // firebase client
		db,
	)

	// Controller
	authController := controller.NewAuth(authService)

	// Routes
	routes.ServeAuth(server, authController, middleware)
}
