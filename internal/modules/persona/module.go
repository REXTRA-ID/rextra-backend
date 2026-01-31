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
	var (
		personaRepository repository.PersonaRepository = repository.NewPersona(db)
		personaService    service.PersonaService       = service.NewPersona(personaRepository, db)
		personaController controller.PersonaController = controller.NewPersona(personaService)
	)

	routes.ServePersona(server, personaController, middleware)

}
