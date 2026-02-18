package controller

import (
	"rextra-backend/internal/modules/token/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	PaymentController interface {
		GetInstructions(ctx *gin.Context)
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
	res, err := c.paymentService.GetInstructions(code)
	if err != nil {
		response.NewFailed("failed get payment instructions", err).Send(ctx)
		return
	}
	response.NewSuccess("success get payment instructions", res).Send(ctx)
}
