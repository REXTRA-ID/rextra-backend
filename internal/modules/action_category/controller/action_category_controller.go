package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/action_category/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	ActionCategoryController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	actionCategoryController struct {
		actionCategoryService service.ActionCategoryService
	}
)

func NewActionCategoryController(service service.ActionCategoryService) ActionCategoryController {
	return &actionCategoryController{
		actionCategoryService: service,
	}
}

func (c *actionCategoryController) Create(ctx *gin.Context) {
	var req dto_request.CreateActionCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.actionCategoryService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create action category", err).Send(ctx)
		return
	}
	response.NewSuccess("action category created successfully", result).Send(ctx)
}

func (c *actionCategoryController) GetAll(ctx *gin.Context) {
	result, err := c.actionCategoryService.GetAll(ctx)
	if err != nil {
		response.NewFailed("failed to retrieve action categories", err).Send(ctx)
		return
	}
	response.NewSuccess("action categories retrieved successfully", result).Send(ctx)
}

func (c *actionCategoryController) GetById(ctx *gin.Context) {
	id := ctx.Param("actionCategoryId")
	result, err := c.actionCategoryService.GetById(ctx, id)
	if err != nil {
		response.NewFailed("failed to retrieve action category", err).Send(ctx)
		return
	}
	response.NewSuccess("action category retrieved successfully", result).Send(ctx)
}

func (c *actionCategoryController) Update(ctx *gin.Context) {
	id := ctx.Param("actionCategoryId")
	var req dto_request.UpdateActionCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.actionCategoryService.Update(ctx, id, req)
	if err != nil {
		response.NewFailed("failed to update action category", err).Send(ctx)
		return
	}
	response.NewSuccess("action category updated successfully", result).Send(ctx)
}

func (c *actionCategoryController) Delete(ctx *gin.Context) {
	id := ctx.Param("actionCategoryId")
	if err := c.actionCategoryService.Delete(ctx, id); err != nil {
		response.NewFailed("failed to delete action category", err).Send(ctx)
		return
	}
	response.NewSuccess("action category deleted successfully", nil).Send(ctx)
}
