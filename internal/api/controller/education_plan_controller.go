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
	EducationPlanController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetEducationPlanById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	educationPlanController struct {
		service service.EducationPlanService
	}
)

func NewEducationPlan(service service.EducationPlanService) EducationPlanController {
	return &educationPlanController{service: service}
}

func (c *educationPlanController) Create(ctx *gin.Context) {
	var req dto_request.CreateEducationPlanRequest
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

	res, err := c.service.Create(ctx.Request.Context(), nil, req)
	
	if err != nil {
		response.NewFailed("failed create education plan", err).Send(ctx)
		return
	}

	response.NewSuccess("success create education plan", res).Send(ctx)
}

func (c *educationPlanController) GetAll(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.service.GetAll(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get all education plan", err).Send(ctx)
		return
	}

	response.NewSuccess("success get all education plan", res).Send(ctx)
}

func (c *educationPlanController) GetEducationPlanById(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	id := ctx.Param("id")
	if id == "" {
		response.NewFailed("failed get education plan id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.service.GetById(ctx.Request.Context(), userId, id)
	if err != nil {
		response.NewFailed("failed get education plan by id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get education plan by id", res).Send(ctx)
}

func (c *educationPlanController) Update(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	id := ctx.Param("id")
	if id == "" {
		response.NewFailed("failed get education plan id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.UpdateEducationPlanRequest
	req.ID = id
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	res, err := c.service.Update(ctx.Request.Context(), userId, req)
	if err != nil {
		response.NewFailed("failed update education plan", err).Send(ctx)
		return
	}

	response.NewSuccess("success update education plan", res).Send(ctx)
}

func (c *educationPlanController) Delete(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	id := ctx.Param("id")
	if id == "" {
		response.NewFailed("failed get education plan id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), userId, id); err != nil {
		response.NewFailed("failed delete education plan", err).Send(ctx)
		return
	}

	response.NewSuccess("success delete education plan", nil).Send(ctx)
}