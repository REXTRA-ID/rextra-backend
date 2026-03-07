package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/feature/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	FeatureController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	featureController struct {
		featureService service.FeatureService
	}
)

func NewFeatureController(featureService service.FeatureService) FeatureController {
	return &featureController{featureService: featureService}
}

func (c *featureController) Create(ctx *gin.Context) {
	var req dto_request.CreateFeatureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.featureService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create feature", err).Send(ctx)
		return
	}
	response.NewSuccess("feature created successfully", result).Send(ctx)
}

func (c *featureController) GetAll(ctx *gin.Context) {
	result, err := c.featureService.GetAll(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve features", err).Send(ctx)
		return
	}
	response.NewSuccess("features retrieved successfully", result).Send(ctx)
}

func (c *featureController) GetById(ctx *gin.Context) {
	id := ctx.Param("featureId")
	result, err := c.featureService.GetById(ctx, id)
	if err != nil {
		response.NewFailed("failed to retrieve feature", err).Send(ctx)
		return
	}
	response.NewSuccess("feature retrieved successfully", result).Send(ctx)
}

func (c *featureController) Update(ctx *gin.Context) {
	id := ctx.Param("featureId")
	var req dto_request.UpdateFeatureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.featureService.Update(ctx, id, req)
	if err != nil {
		response.NewFailed("failed to update feature", err).Send(ctx)
		return
	}
	response.NewSuccess("feature updated successfully", result).Send(ctx)
}

func (c *featureController) Delete(ctx *gin.Context) {
	id := ctx.Param("featureId")
	if err := c.featureService.Delete(ctx, id); err != nil {
		response.NewFailed("failed to delete feature", err).Send(ctx)
		return
	}
	response.NewSuccess("feature deleted successfully", nil).Send(ctx)
}
