package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/action_category/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	ActionCategoryController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
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
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", err).ChangeStatusCode(400).Send(ctx)
		return
	}
	result, err := c.actionCategoryService.Create(ctx, req)
	if err != nil {
		response.NewFailed("Failed to create action category", nil).Send(ctx)
		return
	}
	response.NewSuccess("Action category created successfully", result).Send(ctx)
}

func (c *actionCategoryController) GetAll(ctx *gin.Context) {
	result, err := c.actionCategoryService.GetAll(ctx)
	if err != nil {
		response.NewFailed("Failed get action category", nil).Send(ctx)
		return
	}
	response.NewSuccess("Action category retreived successfully", result).Send(ctx)
}
