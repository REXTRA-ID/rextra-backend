package routes

import (
	"rextra-backend/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func ServeAdmin(app *gin.Engine, admincontroller controller.AdminController) {
	routes := app.Group("/api/v1/admin")
	{
		routes.POST("/vouchers/generate", admincontroller.GenerateVouchers)
	}
}
