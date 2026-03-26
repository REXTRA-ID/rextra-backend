package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/membership/service"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	SubscriptionCycleController interface {
		GetMyCycles(ctx *gin.Context)
		GetAllCycles(ctx *gin.Context)
	}

	subscriptionCycleController struct {
		service service.SubscriptionCycleService
	}
)

func NewSubscriptionCycleController(service service.SubscriptionCycleService) SubscriptionCycleController {
	return &subscriptionCycleController{
		service: service,
	}
}

func (ctrl *subscriptionCycleController) GetMyCycles(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("user id missing", err).ChangeStatusCode(400).Send(ctx)
		return
	}

	cycles, err := ctrl.service.GetByUserID(ctx, userId)
	if err != nil {
		response.NewFailed("failed get my cycles", err).ChangeStatusCode(500).Send(ctx)
		return
	}

	response.NewSuccess("success get my cycles", cycles).Send(ctx)
}

func (ctrl *subscriptionCycleController) GetAllCycles(ctx *gin.Context) {
	var req dto_request.SubscriptionCycleFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("failed get data from body", err).ChangeStatusCode(400).Send(ctx)
		return
	}

	res, err := ctrl.service.GetAllPaginated(ctx, req)
	if err != nil {
		response.NewFailed("failed get all cycles", err).ChangeStatusCode(500).Send(ctx)
		return
	}

	response.NewSuccess("success get all cycles", res).Send(ctx)
}
