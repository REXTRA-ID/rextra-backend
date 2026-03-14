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
		durationService service.PlanDurationService
	}
)

func NewPlanDurationController(svc service.PlanDurationService) PlanDurationController {
	return &planDurationController{durationService: svc}
}

func (c *planDurationController) Create(ctx *gin.Context) {
	var req dto_request.CreatePlanDurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.durationService.Create(ctx, ctx.Param("planId"), req)
	if err != nil { response.NewFailed("failed create", err).Send(ctx); return }
	response.NewSuccess("created", result).Send(ctx)
}

func (c *planDurationController) GetAll(ctx *gin.Context) {
	result, err := c.durationService.GetAllByPlanID(ctx, ctx.Param("planId"))
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *planDurationController) GetById(ctx *gin.Context) {
	result, err := c.durationService.GetById(ctx, ctx.Param("planId"), ctx.Param("durationId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *planDurationController) Update(ctx *gin.Context) {
	var req dto_request.UpdatePlanDurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.durationService.Update(ctx, ctx.Param("planId"), ctx.Param("durationId"), req)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *planDurationController) Delete(ctx *gin.Context) {
	if err := c.durationService.Delete(ctx, ctx.Param("planId"), ctx.Param("durationId")); err != nil { response.NewFailed("failed delete", err).Send(ctx); return }
	response.NewSuccess("deleted", nil).Send(ctx)
}
