package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/kenali_diri/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	KenalidiriAdminController interface {
		GetTestHistory(ctx *gin.Context)
		GetTestDetail(ctx *gin.Context)
		DeleteTestData(ctx *gin.Context)
		ExportTestHistory(ctx *gin.Context)
		GetRiasecCodeList(ctx *gin.Context)
		GetRiasecCodeDetail(ctx *gin.Context)
		UpdateRiasecCode(ctx *gin.Context)
	}

	kenalidiriAdminController struct {
		service service.KenalidiriAdminService
	}
)

func NewKenalidiriAdmin(s service.KenalidiriAdminService) KenalidiriAdminController {
	return &kenalidiriAdminController{service: s}
}

func (c *kenalidiriAdminController) GetTestHistory(ctx *gin.Context) {
	var req dto_request.GetTestHistoryRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetTestHistory(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to fetch test history", err).Send(ctx)
		return
	}

	response.NewSuccess("success get test history", result).Send(ctx)
}

func (c *kenalidiriAdminController) GetTestDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetTestDetail(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to fetch test detail", err).Send(ctx)
		return
	}

	response.NewSuccess("success get test detail", result).Send(ctx)
}

func (c *kenalidiriAdminController) DeleteTestData(ctx *gin.Context) {
	var req dto_request.DeleteTestDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.TestIDs) == 0 {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.service.DeleteTestData(ctx.Request.Context(), req); err != nil {
		response.NewFailed("failed delete test data", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete test data", nil).Send(ctx)
}

func (c *kenalidiriAdminController) ExportTestHistory(ctx *gin.Context) {
	var req dto_request.ExportTestHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.ExportTestHistory(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed export test history", err).Send(ctx)
		return
	}

	response.NewSuccess("success export test history", result).Send(ctx)
}

func (c *kenalidiriAdminController) GetRiasecCodeList(ctx *gin.Context) {
	var req dto_request.GetRiasecCodeListRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetRiasecCodeList(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed get riasec code list", err).Send(ctx)
		return
	}

	response.NewSuccess("success get riasec code list", result).Send(ctx)
}

func (c *kenalidiriAdminController) GetRiasecCodeDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetRiasecCodeDetail(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed get riasec code detail", err).Send(ctx)
		return
	}

	response.NewSuccess("success get riasec code detail", result).Send(ctx)
}

func (c *kenalidiriAdminController) UpdateRiasecCode(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.UpdateRiasecCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.service.UpdateRiasecCode(ctx.Request.Context(), id, req); err != nil {
		response.NewFailed("failed update riasec code", err).Send(ctx)
		return
	}

	response.NewSuccess("success update riasec code", nil).Send(ctx)
}
