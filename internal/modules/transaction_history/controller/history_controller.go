package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/transaction_history/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	HistoryController interface {
		GetTransactions(ctx *gin.Context)
	}

	historyController struct {
		svc service.HistoryService
	}
)

func NewHistoryController(svc service.HistoryService) HistoryController {
	return &historyController{svc: svc}
}

func (c *historyController) GetTransactions(ctx *gin.Context) {
	var filter dto_request.HistoryTransactionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.NewFailed("invalid query params", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	userIDStr, err := utils.GetUserIdFromCtx(ctx)
	if err != nil { response.NewFailed("unauthorized", err).Send(ctx); return }
	
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.GetTransactions(ctx, userID, filter)
	if err != nil {
		response.NewFailed("failed to get transaction history", err).Send(ctx)
		return
	}
	response.NewSuccess("transaction history retrieved", result).Send(ctx)
}
