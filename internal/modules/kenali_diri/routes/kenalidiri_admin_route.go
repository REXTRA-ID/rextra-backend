package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/kenali_diri/controller"

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



		adminRoutes.GET("/riasec-codes", adminController.GetRiasecCodeList)
		adminRoutes.GET("/riasec-codes/:id", adminController.GetRiasecCodeDetail)
		adminRoutes.PUT("/riasec-codes/:id", adminController.UpdateRiasecCode)
	}
}
