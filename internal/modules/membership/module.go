package membership

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/membership/controller"
	"rextra-backend/internal/modules/membership/repository"
	"rextra-backend/internal/modules/membership/routes"
	"rextra-backend/internal/modules/membership/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repositories
	membershipPlanRepository := repository.NewMembershipPlanRepository(db)
	membershipDurationRepository := repository.NewMembershipDurationRepository(db)

	// Services
	membershipPlanService := service.NewMembershipPlanService(membershipPlanRepository, db)
	membershipDurationService := service.NewMembershipDurationService(membershipDurationRepository, db)

	// Controllers
	membershipPlanController := controller.NewMembership(membershipPlanService)
	membershipDurationController := controller.NewMembershipDurationController(membershipDurationService)

	// Routes
	routes.ServeMembershipPlan(server, membershipPlanController, middleware)
	routes.ServeMemberDuration(server, membershipDurationController, middleware)
}
