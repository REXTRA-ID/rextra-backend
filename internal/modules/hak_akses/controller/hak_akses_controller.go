package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/hak_akses/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	HakAksesController interface {
		MakeHakAkses(ctx *gin.Context)
	}

	hakAksesController struct {
		hakAksesService service.HakAksesService
	}
)

func NewHakAksesController(hakAksesService service.HakAksesService) HakAksesController {
	return &hakAksesController{
		hakAksesService: hakAksesService,
	}
}

func (c *hakAksesController) MakeHakAkses(ctx *gin.Context) {

	var data dto_request.MakeHakAksesRequest
	if err := ctx.ShouldBind(&data); err != nil {
		response.NewFailed("Request body is invalid", err, nil).Send(ctx)
		return
	}

	newHakAses, err := c.hakAksesService.MakeHakAkses(ctx, data)
	if err != nil {
		response.NewFailed("Failed to make new hak akses feature", err, nil).Send(ctx)
		return
	}

	response.NewSuccess("Hak akses created successfully", newHakAses).Send(ctx)
}
