package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/membership/controller"

	"github.com/gin-gonic/gin"
)

func ServeMembership(
	app *gin.Engine,
	planController controller.MembershipPlanController,
	durationController controller.PlanDurationController,
	mappingController controller.DurationAccessMappingController,
	cycleController controller.SubscriptionCycleController,
	mw middleware.Middleware,
) {
	// ── Admin routes ─────────────────────────────────────────────────────────
	admin := app.Group("/api/v1/admin/membership")
	admin.Use(mw.Authenticate())
	{
		// MembershipPlan CRUD
		admin.GET("/plan", planController.GetAll)
		admin.GET("/plan/:planId", planController.GetById)
		admin.POST("/plan", planController.Create)
		admin.PUT("/plan/:planId", planController.Update)
		admin.DELETE("/plan/:planId", planController.Delete)

		// PlanDuration — nested di bawah plan
		admin.GET("/plan/:planId/duration", durationController.GetAll)
		admin.GET("/plan/:planId/duration/:durationId", durationController.GetById)
		admin.POST("/plan/:planId/duration", durationController.Create)
		admin.PUT("/plan/:planId/duration/:durationId", durationController.Update)
		admin.DELETE("/plan/:planId/duration/:durationId", durationController.Delete)

		// DurationAccessMapping — nested di bawah duration
		admin.GET("/plan/:planId/duration/:durationId/mapping", mappingController.GetAll)
		admin.GET("/plan/:planId/duration/:durationId/mapping/:mappingId", mappingController.GetById)
		admin.POST("/plan/:planId/duration/:durationId/mapping", mappingController.Create)
		admin.PUT("/plan/:planId/duration/:durationId/mapping/:mappingId", mappingController.Update)
		admin.DELETE("/plan/:planId/duration/:durationId/mapping/:mappingId", mappingController.Delete)

		// Subscription Cycle
		admin.GET("/cycles", cycleController.GetAllCycles)
	}

	// ── User routes ──────────────────────────────────────────────────────────
	user := app.Group("/api/v1/membership")
	user.Use(mw.Authenticate())
	{
		// Katalog plan aktif untuk halaman pilih membership
		user.GET("/catalog", planController.GetCatalog)

		// Riwayat siklus langganan user
		user.GET("/my-cycles", cycleController.GetMyCycles)
	}
}
