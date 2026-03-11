package controller

import (
	"rextra-backend/internal/modules/token/repository"
	"rextra-backend/internal/modules/token/service"
	"rextra-backend/internal/pkg/meta"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	TokenLedgerController interface {
		GetActivity(ctx *gin.Context)
	}
	tokenLedgerController struct {
		tokenLedgerService service.TokenLedgerService
	}
)

func NewTokenLedgerController(tokenLedgerService service.TokenLedgerService) TokenLedgerController {
	return &tokenLedgerController{
		tokenLedgerService: tokenLedgerService,
	}
}

func (c *tokenLedgerController) GetActivity(ctx *gin.Context) {
	m := meta.New(ctx)
	offset, limit := m.GetSkipAndLimit()

	direction := ctx.Query("direction")
	sourceType := ctx.Query("source_type")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	filters := repository.TokenLedgerFilter{}
	if direction != "" {
		filters.Direction = &direction
	}
	if sourceType != "" {
		filters.SourceType = &sourceType
	}
	if startDate != "" {
		filters.StartDate = &startDate
	}
	if endDate != "" {
		filters.EndDate = &endDate
	}

	ledgers, total, err := c.tokenLedgerService.GetActivity(ctx, filters, limit, offset)

	if err != nil {
		response.NewFailed("failed to get all token ledger activities", err).Send(ctx)
		return
	}

	m.Count(int(total))

	response.NewSuccess("success get all token ledger activities", ledgers, m).Send(ctx)

}
