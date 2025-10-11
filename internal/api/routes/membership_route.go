package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeMembership(app *gin.Engine, membershipController controller.MembershipController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/membership")
	{
		routes.POST("")
	}
}
