package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/feature/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	SubFeatureController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	subFeatureController struct {
		subFeatureService service.SubFeatureService
	}
)

func NewSubFeatureController(subFeatureService service.SubFeatureService) SubFeatureController {
	return &subFeatureController{subFeatureService: subFeatureService}
}

func (c *subFeatureController) Create(ctx *gin.Context) {
	featureID := ctx.Param("featureId")
	var req dto_request.CreateSubFeatureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.subFeatureService.Create(ctx, featureID, req)
	if err != nil {
		response.NewFailed("failed to create sub feature", err).Send(ctx)
		return
	}
	response.NewSuccess("sub feature created successfully", result).Send(ctx)
}

func (c *subFeatureController) GetAll(ctx *gin.Context) {
	featureID := ctx.Param("featureId")
	result, err := c.subFeatureService.GetAllByFeatureID(ctx, featureID)
	if err != nil {
		response.NewFailed("failed to retrieve sub features", err).Send(ctx)
		return
	}
	response.NewSuccess("sub features retrieved successfully", result).Send(ctx)
}

func (c *subFeatureController) GetById(ctx *gin.Context) {
	featureID := ctx.Param("featureId")
	id := ctx.Param("subFeatureId")
	result, err := c.subFeatureService.GetById(ctx, featureID, id)
	if err != nil {
		response.NewFailed("failed to retrieve sub feature", err).Send(ctx)
		return
	}
	response.NewSuccess("sub feature retrieved successfully", result).Send(ctx)
}

func (c *subFeatureController) Update(ctx *gin.Context) {
	featureID := ctx.Param("featureId")
	id := ctx.Param("subFeatureId")
	var req dto_request.UpdateSubFeatureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.subFeatureService.Update(ctx, featureID, id, req)
	if err != nil {
		response.NewFailed("failed to update sub feature", err).Send(ctx)
		return
	}
	response.NewSuccess("sub feature updated successfully", result).Send(ctx)
}

func (c *subFeatureController) Delete(ctx *gin.Context) {
	featureID := ctx.Param("featureId")
	id := ctx.Param("subFeatureId")
	if err := c.subFeatureService.Delete(ctx, featureID, id); err != nil {
		response.NewFailed("failed to delete sub feature", err).Send(ctx)
		return
	}
	response.NewSuccess("sub feature deleted successfully", nil).Send(ctx)
}
