package actioncategory

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/action_category/controller"
	"rextra-backend/internal/modules/action_category/repository"
	"rextra-backend/internal/modules/action_category/routes"
	"rextra-backend/internal/modules/action_category/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	actionCategoryRepository := repository.NewActionCategoryRepository(db)

	actionCategoryService := service.NewActionCategoryService(actionCategoryRepository)

	actionCategoryController := controller.NewActionCategoryController(actionCategoryService)

	routes.ServeActionCategory(server, actionCategoryController, middleware)
}
