package user

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/user/controller"
	userRepo "rextra-backend/internal/modules/user/repository"
	"rextra-backend/internal/modules/user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repository (shared with auth module)
	userRepository := userRepo.NewUser(db)

	// Service
	userService := service.NewUser(userRepository, db)

	// Controller
	userController := controller.NewUser(userService)

	// Note: User routes are served by auth module
	_ = userController // Controller created but routes handled by auth module
}
