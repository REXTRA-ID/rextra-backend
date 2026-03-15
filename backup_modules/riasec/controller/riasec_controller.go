package controller

import (
	"rextra-backend/internal/modules/riasec/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	RiasecController interface {
		Create(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
		Update(ctx *gin.Context)
	}

	riasecController struct {
		riasecService service.RiasecService
	}
)

func NewRiasec(riasecService service.RiasecService) RiasecController {
	return &riasecController{
		riasecService: riasecService,
	}
}

func (c *riasecController) Create(ctx *gin.Context) {
	var req dto_request.CreateRiasecRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	req.UserID = userId

	createResult, err := c.riasecService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create riasec", err).Send(ctx)
		return
	}

	response.NewSuccess("success create riasec", createResult).Send(ctx)
}

func (c *riasecController) GetByUserID(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	riasec, err := c.riasecService.GetByUserID(ctx, userId)
	if err != nil {
		response.NewFailed("failed get riasec by id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get riasec by id", riasec).Send(ctx)
}

func (c *riasecController) Update(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.CreateRiasecRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	updateResult, err := c.riasecService.Update(ctx, userId, req)
	if err != nil {
		response.NewFailed("failed update riasec", err).Send(ctx)
		return
	}

	response.NewSuccess("success update riasec", updateResult).Send(ctx)
}
