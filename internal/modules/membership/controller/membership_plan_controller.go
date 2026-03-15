package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	MembershipPlanController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GetCatalog(ctx *gin.Context)
	}

	membershipPlanController struct {
		planService service.MembershipPlanService
	}
)

func NewMembershipPlan(planService service.MembershipPlanService) MembershipPlanController {
	return &membershipPlanController{planService: planService}
}

func (c *membershipPlanController) Create(ctx *gin.Context) {
	var req dto_request.CreateMembershipPlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.planService.Create(ctx, req)
	if err != nil { response.NewFailed("failed create", err).Send(ctx); return }
	response.NewSuccess("created", result).Send(ctx)
}

func (c *membershipPlanController) GetAll(ctx *gin.Context) {
	result, err := c.planService.GetAll(ctx)
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *membershipPlanController) GetById(ctx *gin.Context) {
	result, err := c.planService.GetById(ctx, ctx.Param("planId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *membershipPlanController) Update(ctx *gin.Context) {
	var req dto_request.UpdateMembershipPlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.planService.Update(ctx, ctx.Param("planId"), req)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *membershipPlanController) Delete(ctx *gin.Context) {
	err := c.planService.Delete(ctx, ctx.Param("planId"))
	if err != nil { response.NewFailed("failed delete", err).Send(ctx); return }
	response.NewSuccess("deleted", nil).Send(ctx)
}

func (c *membershipPlanController) GetCatalog(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	result, err := c.planService.GetCatalog(ctx, userID)
	if err != nil { response.NewFailed("failed get catalog", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}
