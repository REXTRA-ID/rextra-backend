package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeEducation(app *gin.Engine, educationController controller.EducationController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/education")
	{
		routes.POST("", middleware.Authenticate(), educationController.Create)
		routes.GET("", middleware.Authenticate(), educationController.GetAll)
	}
}