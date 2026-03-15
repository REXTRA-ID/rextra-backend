package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/user_entitlement/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	UserEntitlementController interface {
		GetMyQuota(ctx *gin.Context)
	}

	userEntitlementController struct {
		svc service.UserEntitlementService
	}
)

func NewUserEntitlementController(svc service.UserEntitlementService) UserEntitlementController {
	return &userEntitlementController{svc: svc}
}

func (c *userEntitlementController) GetMyQuota(ctx *gin.Context) {
	var filter dto_request.MyQuotaFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.NewFailed("failed get query params", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	userIDStr, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", err).Send(ctx); return }
	
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.GetMyQuota(ctx, userID, filter)
	if err != nil {
		response.NewFailed("failed to get entitlement quota", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlement quota retrieved", result).Send(ctx)
}
