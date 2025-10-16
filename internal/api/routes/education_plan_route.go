package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeEducationPlan(app *gin.Engine, educationPlanController controller.EducationPlanController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/education-plan")
	{
		routes.POST("", middleware.Authenticate(), educationPlanController.Create)
		routes.GET("", middleware.Authenticate(), educationPlanController.GetAll)
		routes.GET("/:id", middleware.Authenticate(), educationPlanController.GetEducationPlanById)
		routes.PUT("/:id", middleware.Authenticate(), educationPlanController.Update)
		routes.DELETE("/:id", middleware.Authenticate(), educationPlanController.Delete)
	}
}