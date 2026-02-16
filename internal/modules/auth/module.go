package auth

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/auth/controller"
	auth "rextra-backend/internal/modules/auth/repository"
	"rextra-backend/internal/modules/auth/routes"
	"rextra-backend/internal/modules/auth/service"
	membershipRepo "rextra-backend/internal/modules/membership/repository"
	user "rextra-backend/internal/modules/user/repository"
	mailer "rextra-backend/internal/pkg/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware, mailerService mailer.Mailer) {
	// Repositories
	userRepository := user.NewUser(db)
	sessionRepository := auth.NewSession(db)

	// Membership repositories (cross-module dependency)
	membershipRepository := membershipRepo.NewMembershipRepository(db)
	membershipPlanRepository := membershipRepo.NewMembershipPlanRepository(db)
	membershipDurationRepository := membershipRepo.NewMembershipDurationRepository(db)

	// Service
	authService := service.NewAuth(
		userRepository,
		membershipRepository,
		membershipDurationRepository,
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
