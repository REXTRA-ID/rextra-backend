package riasec

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/riasec/controller"
	"rextra-backend/internal/modules/riasec/repository"
	"rextra-backend/internal/modules/riasec/routes"
	"rextra-backend/internal/modules/riasec/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repository
	riasecRepository := repository.NewRiasec(db)

	// Service
	riasecService := service.NewRiasec(riasecRepository, db)

	// Controller
	riasecController := controller.NewRiasec(riasecService)

	// Routes
	routes.ServeRiasec(server, riasecController, middleware)
}
