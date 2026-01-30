package controller

import (
	"rextra-backend/internal/api/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	CareerRecommendationController interface {
		Create(ctx *gin.Context)
		GetByTestSessionID(ctx *gin.Context)
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

	createResult, err := c.careerRecommendationService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create career recommendation", err).Send(ctx)
		return
	}

	response.NewSuccess("success create career recommendation", createResult).Send(ctx)

}

func (c *careerRecommendationController) GetByTestSessionID(ctx *gin.Context) {
	testSessionIDParam := ctx.Param("id")
	testSessionID, err := strconv.ParseInt(testSessionIDParam, 10, 64)
	if err != nil {
		response.NewFailed("invalid test session id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	careerRecommendation, err := c.careerRecommendationService.GetByTestSessionID(ctx, testSessionID)
	if err != nil {
		response.NewFailed("failed get career recommendation by test session id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get career recommendation by test session id", careerRecommendation).Send(ctx)
}

func (c *careerRecommendationController) Update(ctx *gin.Context) {
	testSessionIDParam := ctx.Param("id")
	testSessionID, err := strconv.ParseInt(testSessionIDParam, 10, 64)
	if err != nil {
		response.NewFailed("invalid test session id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.CreateCareerRecommendationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	updateResult, err := c.careerRecommendationService.Update(ctx, testSessionID, req)
	if err != nil {
		response.NewFailed("failed update career recommendation", err).Send(ctx)
		return
	}

	response.NewSuccess("success update career recommendation", updateResult).Send(ctx)
}
