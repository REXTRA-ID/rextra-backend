package routes

import (
	"rextra-backend/internal/api/kenali_diri/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeKenalidiriAdmin(
	app *gin.Engine,
	adminController controller.KenalidiriAdminController,
	middleware middleware.Middleware,
) {
	adminRoutes := app.Group("/api/v1/admin/kenali-diri")
	adminRoutes.Use(middleware.Authenticate(), middleware.OnlyAdmin())
	{
		adminRoutes.GET("/history", adminController.GetTestHistory)
		adminRoutes.GET("/history/:id", adminController.GetTestDetail)
		adminRoutes.DELETE("/history", adminController.DeleteTestData)
		adminRoutes.POST("/history/export", adminController.ExportTestHistory)

		adminRoutes.GET("/feedback/student", adminController.GetStudentFeedbackList)
		adminRoutes.GET("/feedback/student/stats", adminController.GetStudentFeedbackStats)
		adminRoutes.GET("/feedback/expert", adminController.GetExpertFeedbackList)
		adminRoutes.GET("/feedback/expert/:id", adminController.GetExpertFeedbackDetail)

		adminRoutes.GET("/riasec-codes", adminController.GetRiasecCodeList)
		adminRoutes.GET("/riasec-codes/:id", adminController.GetRiasecCodeDetail)
		adminRoutes.PUT("/riasec-codes/:id", adminController.UpdateRiasecCode)
	}
}
