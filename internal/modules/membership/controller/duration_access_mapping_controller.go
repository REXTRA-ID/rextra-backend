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
	planID := ctx.Param("planId")
	planDurationID := ctx.Param("durationId")
	var req dto_request.CreateDurationAccessMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.mappingService.Create(ctx, planID, planDurationID, req)
	if err != nil {
		response.NewFailed("failed to create access mapping", err).Send(ctx)
		return
	}
	response.NewSuccess("access mapping created successfully", result).Send(ctx)
}

func (c *durationAccessMappingController) GetAll(ctx *gin.Context) {
	planID := ctx.Param("planId")
	planDurationID := ctx.Param("durationId")
	result, err := c.mappingService.GetAllByPlanDurationID(ctx, planID, planDurationID)
	if err != nil {
		response.NewFailed("failed to retrieve access mappings", err).Send(ctx)
		return
	}
	response.NewSuccess("access mappings retrieved successfully", result).Send(ctx)
}

func (c *durationAccessMappingController) GetById(ctx *gin.Context) {
	planID := ctx.Param("planId")
	planDurationID := ctx.Param("durationId")
	id := ctx.Param("mappingId")
	result, err := c.mappingService.GetById(ctx, planID, planDurationID, id)
	if err != nil {
		response.NewFailed("failed to retrieve access mapping", err).Send(ctx)
		return
	}
	response.NewSuccess("access mapping retrieved successfully", result).Send(ctx)
}

func (c *durationAccessMappingController) Update(ctx *gin.Context) {
	planID := ctx.Param("planId")
	planDurationID := ctx.Param("durationId")
	id := ctx.Param("mappingId")
	var req dto_request.UpdateDurationAccessMappingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.mappingService.Update(ctx, planID, planDurationID, id, req)
	if err != nil {
		response.NewFailed("failed to update access mapping", err).Send(ctx)
		return
	}
	response.NewSuccess("access mapping updated successfully", result).Send(ctx)
}

func (c *durationAccessMappingController) Delete(ctx *gin.Context) {
	planID := ctx.Param("planId")
	planDurationID := ctx.Param("durationId")
	id := ctx.Param("mappingId")
	if err := c.mappingService.Delete(ctx, planID, planDurationID, id); err != nil {
		response.NewFailed("failed to delete access mapping", err).Send(ctx)
		return
	}
	response.NewSuccess("access mapping deleted successfully", nil).Send(ctx)
}
