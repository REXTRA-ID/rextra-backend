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
	public := app.Group("/api/v1/membership")
	{
		public.GET("/plan", mw.Authenticate(), planCtrl.GetCatalog)
	}

	admin := app.Group("/api/v1/admin/membership")
	admin.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		// Plans
		plans := admin.Group("/plan")
		{
			plans.POST("", planCtrl.Create)
			plans.GET("", planCtrl.GetAll)
			plans.GET("/:planId", planCtrl.GetById)
			plans.PUT("/:planId", planCtrl.Update)
			plans.DELETE("/:planId", planCtrl.Delete)

			// Durations
			durations := plans.Group("/:planId/duration")
			{
				durations.POST("", durationCtrl.Create)
				durations.GET("", durationCtrl.GetAll)
				durations.GET("/:durationId", durationCtrl.GetById)
				durations.PUT("/:durationId", durationCtrl.Update)
				durations.DELETE("/:durationId", durationCtrl.Delete)

				// Mappings (Benefits)
				mapping := durations.Group("/:durationId/mapping")
				{
					mapping.POST("", mappingCtrl.Create)
					mapping.GET("/:mappingId", mappingCtrl.GetById)
					mapping.PUT("/:mappingId", mappingCtrl.Update)
					mapping.DELETE("/:mappingId", mappingCtrl.Delete)
				}
			}
		}
	}
}
