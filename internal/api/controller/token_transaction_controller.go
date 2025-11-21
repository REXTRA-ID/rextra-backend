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
	TokenTransactionController interface {
		UseToken(ctx *gin.Context)
	}

	tokenTransactionController struct {
		tokenTransactionService service.TokenTransactionService
	}
)

func NewTokenTransactionController(tokenTransactionService service.TokenTransactionService) TokenTransactionController {
	return &tokenTransactionController{
		tokenTransactionService: tokenTransactionService,
	}
}

func (c *tokenTransactionController) UseToken(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.UseTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	useToken, err := c.tokenTransactionService.UseToken(ctx, req, userId)
	if err != nil {
		response.NewFailed("failed to use token", myerror.ProcessingError(err)).Send(ctx)
		return
	}

	response.NewSuccess("success to use token", useToken).Send(ctx)
}
