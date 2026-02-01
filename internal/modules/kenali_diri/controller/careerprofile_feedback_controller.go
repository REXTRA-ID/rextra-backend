package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	dto_req "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/kenali_diri/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
)

type CareerProfileFeedbackController interface {
	GetStudentFeedbacks(ctx *gin.Context)
	GetExpertFeedbacks(ctx *gin.Context)
	GetExpertFeedbackDetail(ctx *gin.Context)
	GetFeedbackMetadata(ctx *gin.Context)
}

type careerProfileFeedbackController struct {
	service service.CareerProfileFeedbackService
}

func NewCareerProfileFeedbackController(s service.CareerProfileFeedbackService) CareerProfileFeedbackController {
	return &careerProfileFeedbackController{service: s}
}

func (c *careerProfileFeedbackController) GetStudentFeedbacks(ctx *gin.Context) {
	var req dto_req.GetStudentFeedbackRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetStudentFeedbacks(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to fetch student feedbacks", err).Send(ctx)
		return
	}

	response.NewSuccess("success get student feedbacks", result).Send(ctx)
}

func (c *careerProfileFeedbackController) GetExpertFeedbacks(ctx *gin.Context) {
	var req dto_req.GetExpertFeedbackRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetExpertFeedbacks(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed to fetch expert feedbacks", err).Send(ctx)
		return
	}

	response.NewSuccess("success get expert feedbacks", result).Send(ctx)
}

func (c *careerProfileFeedbackController) GetExpertFeedbackDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.GetExpertFeedbackDetail(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to fetch expert feedback detail", err).Send(ctx)
		return
	}

	response.NewSuccess("success get expert feedback detail", result).Send(ctx)
}
func (c *careerProfileFeedbackController) GetFeedbackMetadata(ctx *gin.Context) {
	result, err := c.service.GetFeedbackMetadata(ctx.Request.Context())
	if err != nil {
		response.NewFailed("failed to fetch feedback metadata", err).Send(ctx)
		return
	}

	response.NewSuccess("success get feedback metadata", result).Send(ctx)
}
