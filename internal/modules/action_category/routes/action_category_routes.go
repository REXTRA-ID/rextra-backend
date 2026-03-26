package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/action_category/controller"

	"github.com/gin-gonic/gin"
)

func ServeActionCategory(app *gin.Engine, actionCategoryController controller.ActionCategoryController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/action-category")
	routes.Use(middleware.Authenticate())
	{
		routes.GET("", actionCategoryController.GetAll)
		routes.GET("/:actionCategoryId", actionCategoryController.GetById)
		routes.POST("", actionCategoryController.Create)
		routes.PUT("/:actionCategoryId", actionCategoryController.Update)
		routes.DELETE("/:actionCategoryId", actionCategoryController.Delete)
	}
}
