package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/checkout/controller"

	"github.com/gin-gonic/gin"
)

func ServeCheckout(server *gin.Engine, ctrl controller.CheckoutController, mw middleware.Middleware) {
	group := server.Group("/api/v1/checkout")
	group.Use(mw.Authenticate())
	{
		group.GET("/prepare", ctrl.PrepareCheckout)
		group.POST("/calculate", ctrl.CalculatePrice)
		group.POST("/initiate", ctrl.InitiateCheckout)
	}
}
