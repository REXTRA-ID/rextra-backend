package career_recommendation

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/career_recommendation/controller"
	"rextra-backend/internal/modules/career_recommendation/repository"
	"rextra-backend/internal/modules/career_recommendation/routes"
	"rextra-backend/internal/modules/career_recommendation/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	// Repository
	careerRecommendationRepository := repository.NewCareerRecommendation(db)

	// Service
	careerRecommendationService := service.NewCareerRecommendation(careerRecommendationRepository, db)

	// Controller
	careerRecommendationController := controller.NewCareerRecommendation(careerRecommendationService)

	// Routes
	routes.ServeCareerRecommendation(server, careerRecommendationController, middleware)
}
