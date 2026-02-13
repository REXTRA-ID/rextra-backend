package controller

import (
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/token/service"
	"rextra-backend/internal/pkg/meta"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	TopupTransactionController interface {
		GetById(ctx *gin.Context)
		GetAll(ctx *gin.Context)
	}

	topupTransactionController struct {
		topupTransactionService service.TopupTransactionService
	}
)

func NewTopupTransactionController(topupTransactionService service.TopupTransactionService) TopupTransactionController {
	return &topupTransactionController{
		topupTransactionService: topupTransactionService,
	}
}

func (c *topupTransactionController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	topupTransaction, err := c.topupTransactionService.GetByID(ctx, id)
	if err != nil {
		response.NewFailed("failed get topup transaction by id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get topup transaction by id", topupTransaction).Send(ctx)
}

func (c *topupTransactionController) GetAll(ctx *gin.Context) {
	m := meta.New(ctx)
	offset, limit := m.GetSkipAndLimit()

	status := ctx.Query("status")
	transactionType := ctx.Query("type")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	filters := repository.TopupTransactionFilters{}
	if status != "" {
		filters.Status = &status
	}
	if transactionType != "" {
		filters.Type = &transactionType
	}
	if startDate != "" {
		filters.StartDate = &startDate
	}
	if endDate != "" {
		filters.EndDate = &endDate
	}

	transactions, total, err := c.topupTransactionService.GetAll(ctx, filters, limit, offset)
	if err != nil {
		response.NewFailed("failed to get all topup transactions", err).Send(ctx)
		return
	}

	m.Count(int(total))

	response.NewSuccess("success get all topup transactions", transactions, m).Send(ctx)
}
