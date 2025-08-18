package controller

import (
	"rextra-backend/internal/api/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		GetById(ctx *gin.Context)
	}

	userController struct {
		userService service.UserService
	}
)

func NewUser(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) GetById(ctx *gin.Context) {
	userId := ctx.Param("id")
	result, err := c.userService.GetById(ctx.Request.Context(), userId)
	if err != nil {
		response.NewFailed("failed get detail user", err).Send(ctx)
		return
	}

	response.NewSuccess("success get detail user", result).Send(ctx)
}
