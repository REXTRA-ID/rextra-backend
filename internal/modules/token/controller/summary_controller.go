package controller

import (
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/token/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	SummaryController interface {
		GetKPI(ctx *gin.Context)
		GetTrendDirection(ctx *gin.Context)
		GetTrendBySourceType(ctx *gin.Context)
	}

	summaryController struct {
		SummaryService service.TokenSummaryService
	}
)

func NewSummaryController(summaryService service.TokenSummaryService) SummaryController {
	return &summaryController{
		SummaryService: summaryService,
	}
}

func (c *summaryController) GetKPI(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		response.NewFailed("failed get data from query", myerror.New("query start_date or end_date must be provided", myerror.Error_InvalidRequest)).Send(ctx)
		return
	}

	res, err := c.SummaryService.GetKPISummary(ctx, startDate, endDate)
	if err != nil {
		response.NewFailed("failed get kpi summary", err).Send(ctx)
		return
	}
	response.NewSuccess("success get kpi summary", res).Send(ctx)
}

func (c *summaryController) GetTrendDirection(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if startDate == "" || endDate == "" {
		response.NewFailed("failed get data from query", myerror.New("query start_date or end_date must be provided", myerror.Error_InvalidRequest)).Send(ctx)
		return
	}

	res, err := c.SummaryService.GetGraphTrendByDirection(ctx, startDate, endDate)
	if err != nil {
		response.NewFailed("failed get trend summary", err).Send(ctx)
		return
	}
	response.NewSuccess("success get trend summary", res).Send(ctx)
}

func (c *summaryController) GetTrendBySourceType(ctx *gin.Context) {
	sourceType := ctx.Query("source_type")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	if sourceType == "" || startDate == "" || endDate == "" {
		response.NewFailed("failed get data from query", myerror.New("query source_type, start_date or end_date must be provided", myerror.Error_InvalidRequest)).Send(ctx)
		return
	}

	res, err := c.SummaryService.GetGraphTrendBySourceType(ctx, entity.TokenSourceType(sourceType), startDate, endDate)
	if err != nil {
		response.NewFailed("failed get trend summary", err).Send(ctx)
		return
	}
	response.NewSuccess("success get trend summary", res).Send(ctx)
}
