package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	PlanDurationController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	planDurationController struct {
		durationService service.MembershipDurationService
	}
)

func NewPlanDurationController(svc service.MembershipDurationService) PlanDurationController {
	return &planDurationController{
		durationService: svc,
	}
}

func (c *planDurationController) Create(ctx *gin.Context) {

	var req dto_request.CreateMembershipDurationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.durationService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create membership duration", err).Send(ctx)
		return
	}

	response.NewSuccess("membership duration created successfully", result).Send(ctx)
}

func (c *planDurationController) GetAll(ctx *gin.Context) {

	planID := ctx.Param("planId")

	result, err := c.durationService.GetAll(ctx, planID)
	if err != nil {
		response.NewFailed("failed to retrieve membership durations", err).Send(ctx)
		return
	}

	response.NewSuccess("membership durations retrieved successfully", result).Send(ctx)
}

func (c *planDurationController) GetById(ctx *gin.Context) {

	id := ctx.Param("durationId")

	result, err := c.durationService.GetById(ctx, id)
	if err != nil {
		response.NewFailed("failed to retrieve membership duration", err).Send(ctx)
		return
	}

	response.NewSuccess("membership duration retrieved successfully", result).Send(ctx)
}

func (c *planDurationController) Update(ctx *gin.Context) {

	id := ctx.Param("durationId")

	var req dto_request.UpdateMembershipDurationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.durationService.Update(ctx, id, req)
	if err != nil {
		response.NewFailed("failed to update membership duration", err).Send(ctx)
		return
	}

	response.NewSuccess("membership duration updated successfully", result).Send(ctx)
}

func (c *planDurationController) Delete(ctx *gin.Context) {

	id := ctx.Param("durationId")

	if err := c.durationService.Delete(ctx, id); err != nil {
		response.NewFailed("failed to delete membership duration", err).Send(ctx)
		return
	}

	response.NewSuccess("membership duration deleted successfully", nil).Send(ctx)
}
