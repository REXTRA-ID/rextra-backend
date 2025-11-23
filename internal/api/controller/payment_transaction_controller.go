package controller

import (
	"encoding/json"
	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	PaymentTransactionController interface {
		MakeNewTransactionMembership(ctx *gin.Context)
		MakeNewTransactionTokenStandAlone(ctx *gin.Context)
		UpdateTransaction(ctx *gin.Context)
	}

	paymentTransactionController struct {
		paymentTransactionService service.PaymentTransactionService
	}
)

func NewPaymentTransactionController(
	paymentTransactionService service.PaymentTransactionService) PaymentTransactionController {

	return &paymentTransactionController{
		paymentTransactionService: paymentTransactionService,
	}
}

func (c *paymentTransactionController) MakeNewTransactionMembership(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	email, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user email from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.MakeNewTransactionMembershipRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.paymentTransactionService.MakeNewTransactionMembership(ctx, req, userId, email)
	if err != nil {
		response.NewFailed("failed to make new transaction", myerror.ProcessingError(err)).Send(ctx)
		return
	}

	response.NewSuccess("successfully made transaction", result).Send(ctx)
}

func (c *paymentTransactionController) MakeNewTransactionTokenStandAlone(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	email, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user email from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.MakeNewTransactionTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.paymentTransactionService.MakeNewTransactionToken(ctx, req, userId, email)
	if err != nil {
		response.NewFailed("failed to make new transaction", myerror.ProcessingError(err)).Send(ctx)
		return
	}

	response.NewSuccess("successfully made transaction", result).Send(ctx)
}

func (c *paymentTransactionController) UpdateTransaction(ctx *gin.Context) {
	metadata := c.HandleWebook(ctx)

	var req dto_request.PaymentWebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed to get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.paymentTransactionService.UpdateTransaction(ctx, req, metadata)
	if err != nil {
		response.NewFailed("failed to update transaction", err).Send(ctx)
		return
	}

	response.NewSuccess("transaction successfully updated", result).Send(ctx)
}

func (c *paymentTransactionController) HandleWebook(ctx *gin.Context) []byte {
	var payload dto_request.MidtransCallback

	metadata, _ := json.Marshal(payload)

	return metadata
}
