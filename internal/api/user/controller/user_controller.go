package userController

import (
	userService "rextra-backend/internal/api/user/service"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		GetById(ctx *gin.Context)
	}

	userController struct {
		userService userService.UserService
	}
)

func New(userService userService.UserService) UserController {
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
