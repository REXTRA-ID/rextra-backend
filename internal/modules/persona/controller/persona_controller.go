package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/persona/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	PersonaController interface {
		Create(ctx *gin.Context)
		Get(ctx *gin.Context)
		UpdateMission(ctx *gin.Context)
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
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	var req dto_request.CreatePersonaRequest
	req.UserID = userId
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	createResult, err := c.personaService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed create persona", err).Send(ctx)
		return
	}

	response.NewSuccess("success create persona", createResult).Send(ctx)
}

func (c *personaController) Get(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	persona, err := c.personaService.Get(ctx, userId)
	if err != nil {
		response.NewFailed("failed get persona", err).Send(ctx)
		return
	}

	response.NewSuccess("success get persona", persona).Send(ctx)
}

func (c *personaController) UpdateMission(ctx *gin.Context) {
	userId, err := utils.GetUserIdFromCtx(ctx)
	if err != nil {
		response.NewFailed("failed get user id from context", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	var req dto_request.MissionPersonaCompleteRequest
	req.UserID = userId
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("failed get data from body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	persona, err := c.personaService.UpdateMission(ctx, req)

	if err != nil {
		response.NewFailed("failed update persona", err).Send(ctx)
		return
	}

	response.NewSuccess("success update persona mission", persona).Send(ctx)
}
