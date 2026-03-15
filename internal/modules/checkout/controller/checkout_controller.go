package controller

import (
	"errors"
	"net/http"

	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/checkout/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	CheckoutController interface {
		Prepare(ctx *gin.Context)
		Calculate(ctx *gin.Context)
		Initiate(ctx *gin.Context)
		Repeat(ctx *gin.Context)
		Cancel(ctx *gin.Context)
		GetTransaction(ctx *gin.Context)
		GetTransactions(ctx *gin.Context)
		GetPaymentChannels(ctx *gin.Context)
	}

	checkoutController struct {
		svc service.CheckoutService
	}
)

func NewCheckoutController(svc service.CheckoutService) CheckoutController {
	return &checkoutController{svc: svc}
}

func (c *checkoutController) Prepare(ctx *gin.Context) {
	var req dto_request.CheckoutPrepareRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("failed get data from query", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.Prepare(ctx, userID, req)
	if err != nil {
		response.NewFailed("failed to prepare checkout", err).Send(ctx)
		return
	}
	response.NewSuccess("checkout data retrieved", result).Send(ctx)
}

func (c *checkoutController) Calculate(ctx *gin.Context) {
	var req dto_request.CheckoutCalculateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.Calculate(ctx, userID, req)
	if err != nil {
		response.NewFailed("failed to calculate price", err).Send(ctx)
		return
	}
	response.NewSuccess("price calculated", result).Send(ctx)
}

func (c *checkoutController) Initiate(ctx *gin.Context) {
	var req dto_request.CheckoutInitiateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.Initiate(ctx, userID, req)
	if err != nil {
		var pendingErr *service.ErrPendingTransactionExists
		if errors.As(err, &pendingErr) {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "kamu masih punya transaksi yang belum dibayar, selesaikan atau batalkan dulu",
				"data": gin.H{
					"existing_pending_trx_id": pendingErr.TransactionID,
				},
			})
			return
		}
		response.NewFailed("failed to initiate transaction", err).Send(ctx)
		return
	}
	response.NewSuccess("transaction initiated successfully", result).Send(ctx)
}

func (c *checkoutController) Repeat(ctx *gin.Context) {
	transactionID := ctx.Param("transactionId")
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.Repeat(ctx, userID, transactionID)
	if err != nil {
		response.NewFailed("failed to prepare repeat transaction", err).Send(ctx)
		return
	}
	response.NewSuccess("repeat data retrieved", result).Send(ctx)
}

func (c *checkoutController) Cancel(ctx *gin.Context) {
	transactionID := ctx.Param("transactionId")
	var req dto_request.CheckoutCancelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	if err := c.svc.Cancel(ctx, userID, transactionID, req); err != nil {
		response.NewFailed("failed to cancel transaction", err).Send(ctx)
		return
	}
	response.NewSuccess("transaction cancelled successfully", nil).Send(ctx)
}

func (c *checkoutController) GetTransaction(ctx *gin.Context) {
	transactionID := ctx.Param("transactionId")
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	result, err := c.svc.GetTransaction(ctx, userID, transactionID)
	if err != nil {
		response.NewFailed("failed to retrieve transaction", err).Send(ctx)
		return
	}
	response.NewSuccess("transaction retrieved", result).Send(ctx)
}

func (c *checkoutController) GetTransactions(ctx *gin.Context) {
	var req dto_request.MyTransactionFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("failed get data from query", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	userIDStr, _ := utils.GetUserIdFromCtx(ctx)
	userID, _ := uuid.Parse(userIDStr)
	results, total, err := c.svc.GetTransactions(ctx, userID, req)
	if err != nil {
		response.NewFailed("failed to retrieve transactions", err).Send(ctx)
		return
	}
	response.NewSuccess("transactions retrieved", results, gin.H{
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	}).Send(ctx)
}

func (c *checkoutController) GetPaymentChannels(ctx *gin.Context) {
	result, err := c.svc.GetPaymentChannels(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve payment channels", err).Send(ctx)
		return
	}
	response.NewSuccess("payment channels retrieved", result).Send(ctx)
}
