package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/pengaturan/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	PengaturanController interface {
		GetInvoiceSettings(ctx *gin.Context)
		UpdateInvoiceSettings(ctx *gin.Context)
		GetTrxIdSettings(ctx *gin.Context)
		UpdateTrxIdSettings(ctx *gin.Context)
		GetNotifSettings(ctx *gin.Context)
		UpdateNotifSettings(ctx *gin.Context)
	}

	pengaturanController struct {
		service service.PengaturanService
	}
)

func NewPengaturanController(svc service.PengaturanService) PengaturanController {
	return &pengaturanController{service: svc}
}

func (c *pengaturanController) GetInvoiceSettings(ctx *gin.Context) {
	result, err := c.service.GetInvoiceSettings(ctx)
	if err != nil { response.NewFailed("failed get invoice settings", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *pengaturanController) UpdateInvoiceSettings(ctx *gin.Context) {
	adminNameVal, _ := ctx.Get("user_name")
	adminName, _ := adminNameVal.(string)
	var req dto_request.UpdateInvoiceSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid body", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.UpdateInvoiceSettings(ctx, req, adminName)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *pengaturanController) GetTrxIdSettings(ctx *gin.Context) {
	result, err := c.service.GetTrxIdSettings(ctx)
	if err != nil { response.NewFailed("failed get trx settings", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *pengaturanController) UpdateTrxIdSettings(ctx *gin.Context) {
	adminNameVal, _ := ctx.Get("user_name")
	adminName, _ := adminNameVal.(string)
	var req dto_request.UpdateTrxIdSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid body", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.UpdateTrxIdSettings(ctx, req, adminName)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *pengaturanController) GetNotifSettings(ctx *gin.Context) {
	result, err := c.service.GetNotifSettings(ctx)
	if err != nil { response.NewFailed("failed get notif settings", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *pengaturanController) UpdateNotifSettings(ctx *gin.Context) {
	adminNameVal, _ := ctx.Get("user_name")
	adminName, _ := adminNameVal.(string)
	var req dto_request.UpdateNotifSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid body", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.UpdateNotifSettings(ctx, req, adminName)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}
