package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeMemberDuration(r *gin.Engine, membershipDurationController controller.MembershipDurationController, middleware middleware.Middleware) {
	routes := r.Group("/api/v1/membership/duration")
	{
		routes.GET("", membershipDurationController.GetAllMembershipDuration)
	}
}
