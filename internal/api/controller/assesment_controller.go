package controller

import (
	"github.com/gin-gonic/gin"

	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"
)

type (
	AssesmentController interface {
		ValidateHash(ctx *gin.Context)
		GetRiasecQuestion(ctx *gin.Context)
		SubmitRiasecAnswer(ctx *gin.Context)
		GetRiasecResult(ctx *gin.Context)
		GetIkigaiQuestion(ctx *gin.Context)
		SubmitIkigaiAnswer(ctx *gin.Context)
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
	cookieHeader := ctx.Request.Header.Get("Cookie")
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if req.Hash == "" {
		response.NewFailed("failed validate hash", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	setCookies, err := c.assesmentService.ValidateHash(ctx.Request.Context(), req, cookieHeader)
	if err != nil {
		response.NewFailed("failed validate hash", err).Send(ctx)
		return
	}

	for _, cookie := range setCookies {
		ctx.Writer.Header().Add("Set-Cookie", cookie)
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
	userId, err := utils.GetUserIdFromCtx(ctx)
	cookieHeader := ctx.Request.Header.Get("Cookie")
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	var req dto_request.RiasecQuestionSubmitRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if req.Answer == nil {
		response.NewFailed("failed submit riasec answer", myerror.ErrBodyRequest).Send(ctx)
		return
	}


	res, setCookies,err := c.assesmentService.SubmitRiasecAnswer(ctx.Request.Context(), req, userId, cookieHeader)
	
	if err != nil {
		response.NewFailed("failed submit riasec answer", err).Send(ctx)
		return	
	}

	for _, cookie := range setCookies {
		ctx.Writer.Header().Add("Set-Cookie", cookie)
	}

	response.NewSuccess("success submit riasec answer", res).Send(ctx)
}

func (c *assesmentcontroller) GetRiasecResult(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	res, err := c.assesmentService.GetRiasecResult(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get riasec result", err).Send(ctx)
		return
	}

	response.NewSuccess("success get riasec result", res).Send(ctx)
}

func (c *assesmentcontroller) GetIkigaiQuestion(ctx *gin.Context) {
	cookieHeader := ctx.Request.Header.Get("Cookie")
	res, setCookies, err := c.assesmentService.GetIkigaiQuestion(ctx.Request.Context(), cookieHeader)
	
	if err != nil {
		response.NewFailed("failed get ikigai question", err).Send(ctx)
		return
	}


	for _, cookie := range setCookies {
		ctx.Writer.Header().Add("Set-Cookie", cookie)
	}


	response.NewSuccess("success get ikigai question", res).Send(ctx)
}

func (c *assesmentcontroller) SubmitIkigaiAnswer(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	cookieHeader := ctx.Request.Header.Get("Cookie")
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	var req dto_request.IkigaiQuestionSubmitRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if req.Answer == nil {
		response.NewFailed("failed submit ikigai answer", myerror.ErrBodyRequest).Send(ctx)
		return
	}

	res, setCookies, err := c.assesmentService.SubmitIkigaiAnswer(ctx.Request.Context(), userId, cookieHeader, req)
	if err != nil {
		response.NewFailed("failed submit ikigai answer", err).Send(ctx)
		return
	}

	for _, cookie := range setCookies {
		ctx.Writer.Header().Add("Set-Cookie", cookie)
	}

	response.NewSuccess("success submit ikigai answer", res).Send(ctx)
}

	