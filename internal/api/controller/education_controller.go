package controller

import (
	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	EducationController interface {
		Create(ctx *gin.Context)
	}
	educationController struct {
		educationService service.EducationService
	}
)

func NewEducation(educationService service.EducationService) EducationController {
	return &educationController{
		educationService: educationService,
	}
}

func (c *educationController) Create(ctx *gin.Context) {
	var req dto_request.CreateEducationRequest
	
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	req.UserID = userId
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.educationService.CreateEducation(ctx.Request.Context(), req)
	if err != nil {
		response.NewFailed("failed create education", err).Send(ctx)
		return
	}

	response.NewSuccess("success create education", res).Send(ctx)
}
