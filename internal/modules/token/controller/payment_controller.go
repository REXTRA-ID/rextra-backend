package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/token/service"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	PaymentController interface {
		GetInstructions(ctx *gin.Context)
		CreateTransaction(ctx *gin.Context)
		CountPaymentPrice(ctx *gin.Context)
	}

	paymentController struct {
		paymentService service.PaymentService
	}
)

func NewPaymentController(service service.PaymentService) PaymentController {
	return &paymentController{paymentService: service}
}

func (c *paymentController) GetInstructions(ctx *gin.Context) {
	code := ctx.Query("code")
	res, err := c.paymentService.GetInstructions(ctx, code)
	if err != nil {
		response.NewFailed("failed get payment instructions", err).Send(ctx)
		return
	}
	response.NewSuccess("success get payment instructions", res).Send(ctx)
}

func (c *paymentController) CreateTransaction(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", err).Send(ctx)
		return
	}

	var req dto_request.CreateTokenTransactionDTORequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed bind request", err).Send(ctx)
		return
	}

	res, err := c.paymentService.CreateTransaction(ctx, req, userId)
	if err != nil {
		response.NewFailed("failed create transaction", err).Send(ctx)
		return
	}
	response.NewSuccess("success create transaction", res).Send(ctx)
}

func (c *paymentController) CountPaymentPrice(ctx *gin.Context) {
	var req dto_request.CreateTokenTransactionDTORequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed bind request", err).Send(ctx)
		return
	}

	res, err := c.paymentService.CountPaymentPrice(ctx, req)
	if err != nil {
		response.NewFailed("failed count payment price", err).Send(ctx)
		return
	}
	response.NewSuccess("success count payment price", res).Send(ctx)
}
