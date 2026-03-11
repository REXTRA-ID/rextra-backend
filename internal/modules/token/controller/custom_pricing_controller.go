package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/token/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	CustomPricingController interface {
		GetCurrent(ctx *gin.Context)
		CreateNewVersion(ctx *gin.Context)
		GetHistory(ctx *gin.Context)
		ToggleActive(ctx *gin.Context)
	}

	customPricingController struct {
		customPricingService service.CustomPricingService
	}
)

func NewCustomPricingController(customPricingService service.CustomPricingService) CustomPricingController {
	return &customPricingController{
		customPricingService: customPricingService,
	}
}

func (c *customPricingController) GetCurrent(ctx *gin.Context) {
	config, err := c.customPricingService.GetCurrent(ctx)
	if err != nil {
		response.NewFailed("failed get token bundles", err).Send(ctx)
		return
	}
	response.NewSuccess("success get token bundles", config).Send(ctx)
}

func (c *customPricingController) CreateNewVersion(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	var req dto_request.NewCustomPricingDTORequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	err = c.customPricingService.CreateNewConfig(ctx, req, userId)

	if err != nil {
		response.NewFailed("failed create new version", err).Send(ctx)
		return
	}

	response.NewSuccess("success create new version", nil).Send(ctx)
}

func (c *customPricingController) GetHistory(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "10")
	tier := ctx.DefaultQuery("tiers", "false")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		response.NewFailed("failed get limit", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	history, err := c.customPricingService.GetHistory(ctx, limitInt, tier == "true")
	if err != nil {
		response.NewFailed("failed get history", err).Send(ctx)
		return
	}

	response.NewSuccess("success get history", history).Send(ctx)
}

func (c *customPricingController) ToggleActive(ctx *gin.Context) {
	active := ctx.DefaultQuery("active", "false")

	activeBool, err := strconv.ParseBool(active)
	if err != nil {
		response.NewFailed("failed get active", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	err = c.customPricingService.ToggleActive(ctx, activeBool)

	if err != nil {
		response.NewFailed("failed toggle active", err).Send(ctx)
		return
	}

	response.NewSuccess("success toggle active", nil).Send(ctx)
}
