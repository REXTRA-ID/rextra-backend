package controller

import (
	"github.com/gin-gonic/gin"

	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
)

type (
	AssesmentController interface {
		ValidateHash(ctx *gin.Context)
		GetRiasecQuestion(ctx *gin.Context)
		SubmitRiasecAnswer(ctx *gin.Context)
		// GetRiasecResult(ctx *gin.Context)
		// GetIkigaiQuestion(ctx *gin.Context)
		// SubmitIkigaiAnswer(ctx *gin.Context)
		// GetIkigaiResult(ctx *gin.Context)
	}

	assesmentcontroller struct {
		assesmentService service.AssesmentService
	}
)

func NewAssesment(assesmentService service.AssesmentService) AssesmentController {
	return &assesmentcontroller{
		assesmentService: assesmentService,
	}
}

func (c *assesmentcontroller) ValidateHash(ctx *gin.Context) {
	var req dto_request.ValidateHashRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if req.Hash == "" {
		response.NewFailed("failed validate hash", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	if err := c.assesmentService.ValidateHash(ctx.Request.Context(), req); err != nil {
		response.NewFailed("failed validate hash", err).Send(ctx)
		return
	}

	response.NewSuccess("success validate hash", nil).Send(ctx)
}

func (c *assesmentcontroller) GetRiasecQuestion(ctx *gin.Context) {
	res, err := c.assesmentService.GetRiasecQuestion(ctx.Request.Context())
	if err != nil {
		response.NewFailed("failed get riasec question", err).Send(ctx)
		return
	}

	response.NewSuccess("success get riasec question", res).Send(ctx)
}

func (c *assesmentcontroller) SubmitRiasecAnswer(ctx *gin.Context) {
	var req dto_request.RiasecQuestionSubmitRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if req.Answer == nil {
		response.NewFailed("failed submit riasec answer", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	res, err := c.assesmentService.SubmitRiasecAnswer(ctx.Request.Context(), req)
	
	if err != nil {
		response.NewFailed("failed submit riasec answer", err).Send(ctx)
		return	
	}

	response.NewSuccess("success submit riasec answer", res).Send(ctx)
}