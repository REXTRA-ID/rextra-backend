package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/subscription/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	SubscriptionCycleController interface {
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
		GetMyCycles(ctx *gin.Context)
	}

	subscriptionCycleController struct {
		service service.SubscriptionCycleService
	}
)

func NewSubscriptionCycleController(svc service.SubscriptionCycleService) SubscriptionCycleController {
	return &subscriptionCycleController{service: svc}
}

func (c *subscriptionCycleController) GetAll(ctx *gin.Context) {
	var req dto_request.SubscriptionCycleFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("invalid query", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.service.GetAll(ctx, req)
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *subscriptionCycleController) GetById(ctx *gin.Context) {
	result, err := c.service.GetById(ctx, ctx.Param("cycleId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *subscriptionCycleController) GetByUserID(ctx *gin.Context) {
	result, err := c.service.GetByUserID(ctx, ctx.Param("userId"))
	if err != nil { response.NewFailed("failed get for user", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *subscriptionCycleController) GetMyCycles(ctx *gin.Context) {
	userID, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.GetMyCycles(ctx, userID)
	if err != nil { response.NewFailed("failed get mine", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}
