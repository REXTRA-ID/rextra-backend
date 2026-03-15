package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/pengaturan/controller"

	"github.com/gin-gonic/gin"
)

func ServePengaturan(
	app *gin.Engine,
	ctrl controller.PengaturanController,
	mw middleware.Middleware,
) {
	pengaturan := app.Group("/api/v1/pengaturan")
	pengaturan.Use(mw.Authenticate(), mw.OnlyAllow("ADMIN"))
	{
		pengaturan.GET("/invoice", ctrl.GetInvoiceSettings)
		pengaturan.PUT("/invoice", ctrl.UpdateInvoiceSettings)
		pengaturan.GET("/transaction-id", ctrl.GetTrxIdSettings)
		pengaturan.PUT("/transaction-id", ctrl.UpdateTrxIdSettings)
		pengaturan.GET("/notification", ctrl.GetNotifSettings)
		pengaturan.PUT("/notification", ctrl.UpdateNotifSettings)
	}
}
