package kenalidiri

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/kenali_diri/controller"
	"rextra-backend/internal/modules/kenali_diri/repository"
	"rextra-backend/internal/modules/kenali_diri/routes"
	"rextra-backend/internal/modules/kenali_diri/service"
	"rextra-backend/internal/pkg/cache"
	"rextra-backend/internal/pkg/export"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware, cacheService cache.CacheService, exportService export.ExportService) {

	var (
		kenalidiriHistoryRepository    repository.KenalidiriHistoryRepository  = repository.NewKenalidiriHistory(db)
		kenalidiriCategoryRepository   repository.KenalidiriCategoryRepository = repository.NewKenalidiriCategory(db, cacheService)
		kenalidiriRiasecCodeRepository repository.RiasecCodeRepository         = repository.NewRiasecCode(db, cacheService)
		testSessionRepository          repository.TestSessionRepository        = repository.NewTestSession(db)
		kenalidiriRiasecRepository     repository.RiasecRepository             = repository.NewRiasec(db)
		ikigaiRepository               repository.IkigaiRepository             = repository.NewIkigai(db)
		recommendationRepository       repository.RecommendationRepository     = repository.NewRecommendation(db)
		feedbackRepository             repository.FeedbackRepository           = repository.NewFeedback(db)
		careerProfileFeedbackRepo      repository.CareerProfileFeedbackRepository = repository.NewCareerProfileFeedbackRepository(db)
	)

	kenalidiriAdminService := service.NewKenalidiriAdmin(
		kenalidiriHistoryRepository,
		kenalidiriCategoryRepository,
		kenalidiriRiasecCodeRepository,
		testSessionRepository,
		kenalidiriRiasecRepository,
		ikigaiRepository,
		recommendationRepository,
		exportService,
		cacheService,
		db,
	)

	kenalidiriAdminController := controller.NewKenalidiriAdmin(kenalidiriAdminService)
	careerProfileFeedbackService := service.NewCareerProfileFeedbackService(careerProfileFeedbackRepo, feedbackRepository)
	careerProfileFeedbackController := controller.NewCareerProfileFeedbackController(careerProfileFeedbackService)

	routes.ServeKenalidiriAdmin(server, kenalidiriAdminController, middleware)
	routes.ServeCareerProfileFeedback(server, careerProfileFeedbackController, middleware)
}
