package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/subscription/controller"

	"github.com/gin-gonic/gin"
)

func ServeSubscription(
	app *gin.Engine,
	membershipCtrl controller.UserMembershipController,
	cycleCtrl controller.SubscriptionCycleController,
	mw middleware.Middleware,
) {
	admin := app.Group("/api/v1/subscription")
	admin.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		admin.GET("/user", membershipCtrl.GetAll)
		admin.GET("/membership/:membershipId", membershipCtrl.GetById)
		admin.GET("/cycle", cycleCtrl.GetAll)
		admin.GET("/cycle/:cycleId", cycleCtrl.GetById)
		admin.GET("/user/:userId/cycle", cycleCtrl.GetByUserID)
	}

	user := app.Group("/api/v1/my")
	user.Use(mw.Authenticate())
	{
		user.GET("/membership", membershipCtrl.GetMyMembership)
		user.GET("/subscription-cycle", cycleCtrl.GetMyCycles)
		user.POST("/claim-starter", membershipCtrl.ClaimStarter)
	}
}
