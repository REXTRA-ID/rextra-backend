package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/action_category/controller"

	"github.com/gin-gonic/gin"
)

func ServeActionCategory(app *gin.Engine, actionCategoryController controller.ActionCategoryController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/action-category")
	{
		routes.GET("", middleware.Authenticate(), actionCategoryController.GetAll)
		routes.POST("", middleware.Authenticate(), actionCategoryController.Create)
	}
}
