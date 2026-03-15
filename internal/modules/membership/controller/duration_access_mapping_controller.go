package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	DurationAccessMappingController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	durationAccessMappingController struct {
		mappingService service.DurationAccessMappingService
	}
)

func NewDurationAccessMappingController(svc service.DurationAccessMappingService) DurationAccessMappingController {
	return &durationAccessMappingController{mappingService: svc}
}

func (c *durationAccessMappingController) Create(ctx *gin.Context) {
	var req dto_request.CreateDurationAccessMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.mappingService.Create(ctx, ctx.Param("planId"), ctx.Param("durationId"), req)
	if err != nil { response.NewFailed("failed create", err).Send(ctx); return }
	response.NewSuccess("created", result).Send(ctx)
}

func (c *durationAccessMappingController) GetAll(ctx *gin.Context) {
	result, err := c.mappingService.GetAllByPlanDurationID(ctx, ctx.Param("planId"), ctx.Param("durationId"))
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *durationAccessMappingController) GetById(ctx *gin.Context) {
	result, err := c.mappingService.GetById(ctx, ctx.Param("planId"), ctx.Param("durationId"), ctx.Param("mappingId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *durationAccessMappingController) Update(ctx *gin.Context) {
	var req dto_request.UpdateDurationAccessMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.mappingService.Update(ctx, ctx.Param("planId"), ctx.Param("durationId"), ctx.Param("mappingId"), req)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *durationAccessMappingController) Delete(ctx *gin.Context) {
	if err := c.mappingService.Delete(ctx, ctx.Param("planId"), ctx.Param("durationId"), ctx.Param("mappingId")); err != nil { response.NewFailed("failed delete", err).Send(ctx); return }
	response.NewSuccess("deleted", nil).Send(ctx)
}
