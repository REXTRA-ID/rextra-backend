package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeMembershipPlan(app *gin.Engine, membershipPlanController controller.MembershipPlanController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/membership/plan")
	{
		routes.GET("", membershipPlanController.GetAllMembershipPlan)
	}
}
