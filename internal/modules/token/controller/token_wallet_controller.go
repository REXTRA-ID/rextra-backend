package controller

import (
	"rextra-backend/internal/modules/token/service"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	TokenWalletController interface {
		// USER
		GetMyBalance(ctx *gin.Context)
		GetMyHistory(ctx *gin.Context)
		GetMyHistoryDetail(ctx *gin.Context)
	}

	tokenWalletController struct {
		walletService service.WalletService
	}
)

func NewTokenWalletController(tokenWalletService service.WalletService) TokenWalletController {
	return &tokenWalletController{
		walletService: tokenWalletService,
	}
}

func (c *tokenWalletController) GetMyBalance(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	res, err := c.walletService.GetMyBalance(ctx, userId)
	if err != nil {
		response.NewFailed("failed get balance", err).Send(ctx)
		return
	}

	response.NewSuccess("success get me", res).Send(ctx)
}

func (c *tokenWalletController) GetMyHistory(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)

	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	if page < 1 || size < 1 {
		response.NewFailed("invalid page or size", err).Send(ctx)
		return
	}

	res, resPage, err := c.walletService.GetMyHistory(ctx, userId, page, size)
	if err != nil {
		response.NewFailed("failed get history", err).Send(ctx)
		return
	}

	response.NewSuccess("success get history", res, resPage).Send(ctx)
}

func (c *tokenWalletController) GetMyHistoryDetail(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)

	if err != nil {
		response.NewFailed("failed get user id", err).Send(ctx)
		return
	}

	id := ctx.Param("id")

	res, err := c.walletService.GetMyHistoryDetail(ctx, userId, id)
	if err != nil {
		response.NewFailed("failed get history detail", err).Send(ctx)
		return
	}

	response.NewSuccess("success get history detail", res).Send(ctx)
}
