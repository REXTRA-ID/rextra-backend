package subscription

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/subscription/controller"
	"rextra-backend/internal/modules/subscription/repository"
	"rextra-backend/internal/modules/subscription/routes"
	"rextra-backend/internal/modules/subscription/service"

	membershipRepo "rextra-backend/internal/modules/membership/repository"
	userRepo "rextra-backend/internal/modules/user/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	membershipRepository := membershipRepo.NewUserMembershipRepository(db)
	cycleRepository := repository.NewSubscriptionCycleRepository(db)
	userRepository := userRepo.NewUserRepository(db)

	membershipService := service.NewUserMembershipService(membershipRepository, userRepository)
	cycleService := service.NewSubscriptionCycleService(cycleRepository, membershipRepository)

	membershipController := controller.NewUserMembershipController(membershipService)
	cycleController := controller.NewSubscriptionCycleController(cycleService)

	routes.ServeSubscription(server, membershipController, cycleController, middleware)
}
