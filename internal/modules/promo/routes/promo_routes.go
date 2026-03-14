package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/promo/controller"

	"github.com/gin-gonic/gin"
)

func ServePromo(
	app *gin.Engine,
	ctrl controller.DiscountController,
	mw middleware.Middleware,
) {
	base := app.Group("/api/v1/promo")
	base.Use(mw.Authenticate())

	base.POST("/validate", ctrl.ValidateCode)

	adminOnly := base.Group("")
	adminOnly.Use(mw.OnlyAllow("ADMIN"))
	{
		adminOnly.GET("", ctrl.GetAll)
		adminOnly.POST("", ctrl.Create)
		adminOnly.GET("/:discountId", ctrl.GetByID)
		adminOnly.PUT("/:discountId", ctrl.Update)
		adminOnly.DELETE("/:discountId", ctrl.Delete)
		adminOnly.GET("/:discountId/redemption", ctrl.GetRedemptions)
	}
}
