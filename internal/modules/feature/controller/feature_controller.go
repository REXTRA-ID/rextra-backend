package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/feature/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	FeatureController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
	}

	featureController struct {
		featureService service.FeatureService
	}
)

func NewFeatureController(service service.FeatureService) FeatureController {
	return &featureController{
		featureService: service,
	}
}

func (c *featureController) Create(ctx *gin.Context) {
	var req dto_request.CreateFeatureRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).ChangeStatusCode(400).Send(ctx)
		return
	}
	result, err := c.featureService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create new feature", err).Send(ctx)
		return
	}
	response.NewSuccess("Feature created successfully", result).Send(ctx)
}

func (c *featureController) GetAll(ctx *gin.Context) {
	result, err := c.featureService.GetAll(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve feature", err).Send(ctx)
		return
	}
	response.NewSuccess("Feature retrieved successfully", result).Send(ctx)
}

func (c *featureController) GetById(ctx *gin.Context) {
	featureId := ctx.Param("featureId")
	result, err := c.featureService.GetById(ctx, featureId)
	if err != nil {
		response.NewFailed("failed to retrieve feature", err).Send(ctx)
		return
	}
	response.NewSuccess("Feature retrieved successfully", result).Send(ctx)
}
