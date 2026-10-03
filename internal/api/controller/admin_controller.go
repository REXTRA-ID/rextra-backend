package controller

import (
	"github.com/gin-gonic/gin"
	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
)

type (
	AdminController interface {
		GenerateVouchers(ctx *gin.Context)
	}

	admincontroller struct {
		adminService service.AdminService
	}
)

func NewAdmin(adminService service.AdminService) AdminController {
	return &admincontroller{
		adminService: adminService,
	}
}

func (c *admincontroller) GenerateVouchers(ctx *gin.Context) {
	var req dto_request.GenerateVouchersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.adminService.GenerateVouchers(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed generate vouchers", err).Send(ctx)
		return
	}

	response.NewSuccess("success generate vouchers", res).Send(ctx)
}
