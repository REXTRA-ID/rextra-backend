package controller

import (
	"rextra-backend/internal/modules/career_recommendation/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	CareerRecommendationController interface {
		Create(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
		Update(ctx *gin.Context)
	}

	careerRecommendationController struct {
		careerRecommendationService service.CareerRecommendationService
	}
)

func NewCareerRecommendation(careerRecommendationService service.CareerRecommendationService) CareerRecommendationController {
	return &careerRecommendationController{
		careerRecommendationService: careerRecommendationService,
	}
}

func (c *careerRecommendationController) Create(ctx *gin.Context) {
	var req dto_request.CreateCareerRecommendationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	req.UserID = userId

	createResult, err := c.careerRecommendationService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create career recommendation", err).Send(ctx)
		return
	}

	response.NewSuccess("success create career recommendation", createResult).Send(ctx)
}

func (c *careerRecommendationController) GetByUserID(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	careerRecommendation, err := c.careerRecommendationService.GetByUserID(ctx, userId)
	if err != nil {
		response.NewFailed("failed get career recommendation by id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get career recommendation by id", careerRecommendation).Send(ctx)
}

func (c *careerRecommendationController) Update(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.CreateCareerRecommendationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	updateResult, err := c.careerRecommendationService.Update(ctx, userId, req)
	if err != nil {
		response.NewFailed("failed update career recommendation", err).Send(ctx)
		return
	}

	response.NewSuccess("success update career recommendation", updateResult).Send(ctx)
}
