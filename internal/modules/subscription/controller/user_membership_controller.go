package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/subscription/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	UserMembershipController interface {
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		GetMyMembership(ctx *gin.Context)
		ClaimStarter(ctx *gin.Context)
		GetDashboard(ctx *gin.Context)
	}

	userMembershipController struct {
		service service.UserMembershipService
	}
)

func NewUserMembershipController(svc service.UserMembershipService) UserMembershipController {
	return &userMembershipController{service: svc}
}

func (c *userMembershipController) GetAll(ctx *gin.Context) {
	var req dto_request.UserMembershipFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("invalid query", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.service.GetAll(ctx, req)
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *userMembershipController) GetById(ctx *gin.Context) {
	result, err := c.service.GetById(ctx, ctx.Param("membershipId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *userMembershipController) GetMyMembership(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.GetMyMembership(ctx, userID)
	if err != nil { response.NewFailed("failed get mine", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *userMembershipController) ClaimStarter(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", myerror.InvalidRequest(err)).Send(ctx); return }
	if err := c.service.ClaimStarter(ctx, userID); err != nil {
		response.NewFailed("failed to claim starter plan", err).Send(ctx)
		return
	}
	response.NewSuccess("starter plan claimed successfully", nil).Send(ctx)
}

func (c *userMembershipController) GetDashboard(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.GetDashboard(ctx, userID)
	if err != nil { response.NewFailed("failed get dashboard", err).Send(ctx); return }
	response.NewSuccess("dashboard retrieved", result).Send(ctx)
}
