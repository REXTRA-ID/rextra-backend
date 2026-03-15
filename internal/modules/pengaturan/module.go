package pengaturan

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/pengaturan/controller"
	"rextra-backend/internal/modules/pengaturan/repository"
	"rextra-backend/internal/modules/pengaturan/routes"
	"rextra-backend/internal/modules/pengaturan/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, mw middleware.Middleware) {
	repo := repository.NewPengaturanRepository(db)
	svc := service.NewPengaturanService(repo)
	ctrl := controller.NewPengaturanController(svc)
	routes.ServePengaturan(server, ctrl, mw)
}
