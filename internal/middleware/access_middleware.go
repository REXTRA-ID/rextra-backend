package middleware

import (
	"fmt"
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func (m AccessFeatureMiddleware) AccessFeature(ctx *gin.Context, feature, action string) {
	userMember, exist := ctx.Get("membership")
	if !exist {
		response.NewFailed("Cannot access feature for non-member", nil, nil)
		return
	}

	accessData := dto_request.CheckAccessRequest{
		MembershipPlan: userMember.(string),
		Feature:        feature,
		Action:         action,
	}

	access, err := m.HakAksesService.CheckAccess(ctx, accessData)
	if err != nil {
		response.NewFailed("Cannot access feature", err, nil).Send(ctx)
		return
	}

	if access {
		ctx.Next()
	}

	msg := fmt.Sprintf("This feature is not available on your %s plan", userMember.(string))

	response.NewFailed(msg, err, nil).Send(ctx)
}
