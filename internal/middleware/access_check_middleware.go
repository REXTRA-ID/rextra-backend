package middleware

import (
	"fmt"

	accessCheckService "rextra-backend/internal/modules/access_check/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m Middleware) AccessCheck(entitlementKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr, exists := ctx.Get("user_id")
		if !exists {
			res := response.NewFailed("Unauthorized", myerror.New("user ID not found in context", myerror.Error_Unauthorized))
			res.SendWithAbort(ctx)
			return
		}

		userID, err := uuid.Parse(fmt.Sprintf("%v", userIDStr))
		if err != nil {
			res := response.NewFailed("Invalid User ID", myerror.New("invalid user ID format", myerror.Error_Unauthorized))
			res.SendWithAbort(ctx)
			return
		}

		if m.accessCheckSvc == nil {
			// Fail-safe in case accessCheckSvc is not injected
			res := response.NewFailed("Internal Server Error", myerror.New("access check service not configured", myerror.SystemError))
			res.SendWithAbort(ctx)
			return
		}

		req := accessCheckService.CheckAccessRequest{
			UserID:         userID,
			EntitlementKey: entitlementKey,
		}

		resp, err := m.accessCheckSvc.CheckAccess(ctx.Request.Context(), req)
		if err != nil {
			res := response.NewFailed("Access Check Failed", err)
			res.SendWithAbort(ctx)
			return
		}

		if !resp.Granted {
			// Return forbidden with the human-readable deny reason
			res := response.NewFailed(resp.DenyReason, myerror.New(fmt.Sprintf("restriction_type: %s", resp.RestrictionType), myerror.Error_Unauthorized))
			res.SendWithAbort(ctx)
			return
		}

		if resp.QuotaRemaining != nil {
			ctx.Set("quota_remaining", *resp.QuotaRemaining)
		}

		ctx.Next()
	}
}
