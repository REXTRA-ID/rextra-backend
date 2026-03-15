package controller

import (
	"rextra-backend/internal/modules/persona/service"
	dto_request "rextra-backend/internal/dto/request"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	PersonaController interface {
		Create(ctx *gin.Context)
		GetByUserID(ctx *gin.Context)
		Update(ctx *gin.Context)
	}

	personaController struct {
		personaService service.PersonaService
	}
)

func NewPersona(personaService service.PersonaService) PersonaController {
	return &personaController{
		personaService: personaService,
	}
}

func (c *personaController) Create(ctx *gin.Context) {
	var req dto_request.CreatePersonaRequest
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

	createResult, err := c.personaService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create persona", err).Send(ctx)
		return
	}

	response.NewSuccess("success create persona", createResult).Send(ctx)
}

func (c *personaController) GetByUserID(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	persona, err := c.personaService.GetByUserID(ctx, userId)
	if err != nil {
		response.NewFailed("failed get persona by id", err).Send(ctx)
		return
	}

	response.NewSuccess("success get persona by id", persona).Send(ctx)
}

func (c *personaController) Update(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.CreatePersonaRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	updateResult, err := c.personaService.Update(ctx, userId, req)
	if err != nil {
		response.NewFailed("failed update persona", err).Send(ctx)
		return
	}

	response.NewSuccess("success update persona", updateResult).Send(ctx)
}
