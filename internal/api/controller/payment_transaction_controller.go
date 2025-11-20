package controller

import (
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
		UpdateTransactionFromWebhook(ctx *gin.Context)
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

	var req dto_request.MakeNewTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.paymentTransactionService.MakeNewTransaction(ctx, req, userId)
	if err != nil {
		response.NewFailed("failed to make new transaction", myerror.ProcessingError(err)).Send(ctx)
		return
	}

	response.NewSuccess("successfully made transaction", result).Send(ctx)
}

func (c *paymentTransactionController) UpdateTransactionFromWebhook(ctx *gin.Context) {

}

func (c *paymentTransactionController) MakeNewTransactionTokenStandAlone(ctx *gin.Context) {

}
