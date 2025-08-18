package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeCareerRecommendation(app *gin.Engine, careerRecommendationcontroller controller.CareerRecommendationController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/user/:id/career-recommendation")
	{
		routes.POST("", middleware.Authenticate(), careerRecommendationcontroller.Create)
		routes.GET("", middleware.Authenticate(), careerRecommendationcontroller.GetByUserID)
		routes.PUT("", middleware.Authenticate(), careerRecommendationcontroller.Update)
	}
}
