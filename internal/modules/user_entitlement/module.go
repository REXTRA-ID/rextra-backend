package user_entitlement

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/user_entitlement/controller"
	"rextra-backend/internal/modules/user_entitlement/repository"
	"rextra-backend/internal/modules/user_entitlement/routes"
	"rextra-backend/internal/modules/user_entitlement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	repo := repository.NewUserEntitlementRepository(db)
	svc := service.NewUserEntitlementService(repo)
	ctrl := controller.NewUserEntitlementController(svc)
	routes.ServeUserEntitlement(server, ctrl, middleware)
}
