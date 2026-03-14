package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/membership/controller"

	"github.com/gin-gonic/gin"
)

func ServeMembership(
	app *gin.Engine,
	planCtrl controller.MembershipPlanController,
	durationCtrl controller.PlanDurationController,
	mappingCtrl controller.DurationAccessMappingController,
	mw middleware.Middleware,
) {
	// User Routes
	userRoutes := app.Group("/api/v1/membership")
	{
		userRoutes.GET("/plan", planCtrl.GetCatalog)
	}

	// Admin Routes
	adminRoutes := app.Group("/api/v1/admin/membership")
	adminRoutes.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		// Plan CRUD
		adminRoutes.GET("/plan", planCtrl.GetAll)
		adminRoutes.POST("/plan", planCtrl.Create)
		adminRoutes.GET("/plan/:planId", planCtrl.GetById)
		adminRoutes.PUT("/plan/:planId", planCtrl.Update)
		adminRoutes.DELETE("/plan/:planId", planCtrl.Delete)

		// Duration CRUD (Nested)
		duration := adminRoutes.Group("/plan/:planId/duration")
		{
			duration.GET("", durationCtrl.GetAll)
			duration.POST("", durationCtrl.Create)
			duration.GET("/:durationId", durationCtrl.GetById)
			duration.PUT("/:durationId", durationCtrl.Update)
			duration.DELETE("/:durationId", durationCtrl.Delete)

			// Access Mapping CRUD (Nested)
			mapping := duration.Group("/:durationId/mapping")
			{
				mapping.GET("", mappingCtrl.GetAll)
				mapping.POST("", mappingCtrl.Create)
				mapping.GET("/:mappingId", mappingCtrl.GetById)
				mapping.PUT("/:mappingId", mappingCtrl.Update)
				mapping.DELETE("/:mappingId", mappingCtrl.Delete)
			}
		}
	}
}
