package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	MembershipPlanController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GetCatalog(ctx *gin.Context) // user endpoint
	}

	membershipPlanController struct {
		planService service.MembershipPlanService
	}
)

func NewMembershipPlanController(svc service.MembershipPlanService) MembershipPlanController {
	return &membershipPlanController{planService: svc}
}

func (c *membershipPlanController) Create(ctx *gin.Context) {
	var req dto_request.CreateMembershipPlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.planService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create membership plan", err).Send(ctx)
		return
	}
	response.NewSuccess("membership plan created successfully", result).Send(ctx)
}

func (c *membershipPlanController) GetAll(ctx *gin.Context) {
	result, err := c.planService.GetAll(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve membership plans", err).Send(ctx)
		return
	}
	response.NewSuccess("membership plans retrieved successfully", result).Send(ctx)
}

func (c *membershipPlanController) GetById(ctx *gin.Context) {
	id := ctx.Param("planId")
	result, err := c.planService.GetById(ctx, id)
	if err != nil {
		response.NewFailed("failed to retrieve membership plan", err).Send(ctx)
		return
	}
	response.NewSuccess("membership plan retrieved successfully", result).Send(ctx)
}

func (c *membershipPlanController) Update(ctx *gin.Context) {
	id := ctx.Param("planId")
	var req dto_request.UpdateMembershipPlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.planService.Update(ctx, id, req)
	if err != nil {
		response.NewFailed("failed to update membership plan", err).Send(ctx)
		return
	}
	response.NewSuccess("membership plan updated successfully", result).Send(ctx)
}

func (c *membershipPlanController) Delete(ctx *gin.Context) {
	id := ctx.Param("planId")
	if err := c.planService.Delete(ctx, id); err != nil {
		response.NewFailed("failed to delete membership plan", err).Send(ctx)
		return
	}
	response.NewSuccess("membership plan deleted successfully", nil).Send(ctx)
}

func (c *membershipPlanController) GetCatalog(ctx *gin.Context) {
	result, err := c.planService.GetCatalog(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve membership catalog", err).Send(ctx)
		return
	}
	response.NewSuccess("membership catalog retrieved successfully", result).Send(ctx)
}
