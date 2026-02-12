package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/token/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	TokenBundleController interface {
		// FOR ALL ROLE
		GetAll(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		// ADMIN ONLY
		Create(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	tokenBundleController struct {
		bundleService service.BundleService
	}
)

func NewTokenBundleController(tokenBundleService service.BundleService) TokenBundleController {
	return &tokenBundleController{
		bundleService: tokenBundleService,
	}
}

func (c *tokenBundleController) GetAll(ctx *gin.Context) {
	bundles, err := c.bundleService.GetAll(ctx)
	if err != nil {
		response.NewFailed("failed get token bundles", err).Send(ctx)
		return
	}

	response.NewSuccess("token bundles retrieved successfully", bundles).Send(ctx)
}

func (c *tokenBundleController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	bundle, err := c.bundleService.GetByID(ctx, id)
	if err != nil {
		response.NewFailed("failed get token bundle", err).Send(ctx)
		return
	}

	response.NewSuccess("token bundle retrieved successfully", bundle).Send(ctx)
}
func (c *tokenBundleController) Create(ctx *gin.Context) {
	var req dto_request.TokenBundleDTORequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.bundleService.Create(ctx, req)

	if err != nil {
		response.NewFailed("failed get token bundles", err).Send(ctx)
		return
	}

	response.NewSuccess("token bundles retrieved successfully", res).Send(ctx)
}

func (c *tokenBundleController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req dto_request.TokenBundleDTORequest

	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.bundleService.Update(ctx, id, req)

	if err != nil {
		response.NewFailed("failed get token bundles", err).Send(ctx)
		return
	}

	response.NewSuccess("token bundles updated successfully", res).Send(ctx)
}

func (c *tokenBundleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.bundleService.Delete(ctx, id)
	if err != nil {
		response.NewFailed("failed to delete token bundle", err).Send(ctx)
		return
	}

	response.NewSuccess("token bundle deleted successfully", nil).Send(ctx)
}
