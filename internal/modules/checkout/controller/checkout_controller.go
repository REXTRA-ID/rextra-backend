package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/checkout/service"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	CheckoutController interface {
		PrepareCheckout(ctx *gin.Context)
		CalculatePrice(ctx *gin.Context)
		InitiateCheckout(ctx *gin.Context)
	}

	checkoutController struct {
		service service.CheckoutService
	}
)

func NewCheckoutController(service service.CheckoutService) CheckoutController {
	return &checkoutController{service}
}

func (c *checkoutController) PrepareCheckout(ctx *gin.Context) {
	var req dto_request.CheckoutPrepareRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.NewFailed("Invalid request", err).Send(ctx)
		return
	}

	userID, _ := utils.GetUserIdFromCtx(ctx)
	res, err := c.service.PrepareCheckout(ctx.Request.Context(), req, userID)
	if err != nil {
		response.NewFailed("Failed to prepare checkout", err).Send(ctx)
		return
	}

	response.NewSuccess("Success", res).Send(ctx)
}

func (c *checkoutController) CalculatePrice(ctx *gin.Context) {
	var req dto_request.CheckoutCalculateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("Invalid request", err).Send(ctx)
		return
	}

	userID, _ := utils.GetUserIdFromCtx(ctx)
	res, err := c.service.CalculatePrice(ctx.Request.Context(), req, userID)
	if err != nil {
		response.NewFailed("Failed to calculate price", err).Send(ctx)
		return
	}

	response.NewSuccess("Success", res).Send(ctx)
}

func (c *checkoutController) InitiateCheckout(ctx *gin.Context) {
	var req dto_request.CheckoutInitiateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("Invalid request", err).Send(ctx)
		return
	}

	userID, _ := utils.GetUserIdFromCtx(ctx)
	res, err := c.service.InitiateCheckout(ctx.Request.Context(), req, userID)
	if err != nil {
		response.NewFailed("Failed to initiate checkout", err).Send(ctx)
		return
	}

	response.NewSuccess("Success", res).Send(ctx)
}
