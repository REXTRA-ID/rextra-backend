package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/user_entitlement/controller"

	"github.com/gin-gonic/gin"
)

func ServeUserEntitlement(server *gin.Engine, ctrl controller.UserEntitlementController, m middleware.Middleware) {
	r := server.Group("/api/v1/entitlement", m.Authenticate())
	{
		r.GET("/my-quota", ctrl.GetMyQuota)
	}
}
