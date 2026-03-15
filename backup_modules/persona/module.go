package persona

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/persona/controller"
	"rextra-backend/internal/modules/persona/repository"
	"rextra-backend/internal/modules/persona/routes"
	"rextra-backend/internal/modules/persona/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repository
	personaRepository := repository.NewPersona(db)

	// Service
	personaService := service.NewPersona(personaRepository, db)

	// Controller
	personaController := controller.NewPersona(personaService)

	// Routes
	routes.ServePersona(server, personaController, middleware)
}
